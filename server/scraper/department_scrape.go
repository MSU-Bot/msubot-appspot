package scraper

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/SpencerCornish/msubot-appspot/server/models"
	"github.com/SpencerCornish/msubot-appspot/server/serverutils"
	log "github.com/sirupsen/logrus"
)

type resultPayload struct {
	Dept         string
	Error        error
	coursesFound int
	execTimeMs   int64
}

type job struct {
	Department string
	Term       string
}

type clients struct {
	HttpClient *http.Client
	fbClient   *firestore.Client
}

type course struct {
	Number        string
	Name          string
	TotalSections int
	TotalSeats    int
}

const springPostfix int = 30
const summerPostfix int = 50
const fallPostfix int = 70

const numConcurrentWorkers int = 20

func getTerms() []string {
	now := time.Now()
	terms := make([]string, 3)

	// If it's before or during april, use this year, otherwise use next year
	if now.Month() <= 4 {
		terms[0] = fmt.Sprintf("%d%d", now.Year(), springPostfix)
	} else {
		terms[0] = fmt.Sprintf("%d%d", now.Year()+1, springPostfix)
	}

	// If it's before or during august, use this year, otherwise use next year
	if now.Month() <= 8 {
		terms[1] = fmt.Sprintf("%d%d", now.Year(), summerPostfix)
	} else {
		terms[1] = fmt.Sprintf("%d%d", now.Year()+1, summerPostfix)
	}

	// If it's before or during october, use this year, otherwise use next year
	if now.Month() <= 10 {
		terms[2] = fmt.Sprintf("%d%d", now.Year(), fallPostfix)
	} else {
		terms[2] = fmt.Sprintf("%d%d", now.Year()+1, fallPostfix)
	}
	return terms
}

func HandleDepartmentRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	log.WithContext(ctx).Infof("Context loaded. Starting execution.")

	fbClient := serverutils.GetFirebaseClient(ctx)
	departmentBatch := fbClient.Batch()

	clientBundle := &clients{
		HttpClient: client,
		fbClient:   fbClient,
	}

	depts, err := fetchAndUpdateDepartments(ctx, departmentBatch, clientBundle)
	if err != nil {
		errorStr := fmt.Sprintf("Department Scrape Failed with error: %v", err)
		http.Error(w, errorStr, http.StatusInternalServerError)
		return
	}

	_, err = departmentBatch.Commit(ctx)
	if err != nil {
		errorStr := fmt.Sprintf("Failed to commit batch: %v", err)
		http.Error(w, errorStr, http.StatusInternalServerError)
		return
	}

	terms := getTerms()

	for _, term := range terms {
		err = fetchCoursesForTerm(ctx, depts, clientBundle, term)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(200)
}

func fetchCoursesForTerm(ctx context.Context, depts []*models.Department, clientBundle *clients, term string) error {

	jobs := make(chan *job, len(depts))
	results := make(chan *resultPayload, len(depts))

	for id := 1; id <= numConcurrentWorkers; id++ {
		go worker(ctx, id, clientBundle, jobs, results)
	}

	for _, dept := range depts {
		jobs <- &job{
			Department: dept.Id,
			Term:       term,
		}
	}

	totalJobs := len(jobs)
	log.Info("Enqueued %d jobs", totalJobs)
	close(jobs)

	for i := 1; i <= totalJobs; i++ {
		result := <-results
		log.Infof("Job #%d/%d Complete! %v", i, totalJobs, result)
	}

	log.Info("Done waiting for results")

	return nil
}

func worker(ctx context.Context, id int, clients *clients, jobs <-chan *job, result chan<- *resultPayload) {
	for job := range jobs {
		batch := clients.fbClient.Batch()

		start := time.Now()
		log.Infof("Requesting courses for dept:%s Term:%s", job.Department, job.Term)

		resp, err := serverutils.MakeAtlasSectionRequest(clients.HttpClient, job.Term, job.Department, "")
		if err != nil {
			result <- &resultPayload{
				Dept:         job.Department,
				Error:        err,
				coursesFound: 0,
				execTimeMs:   time.Since(start).Milliseconds(),
			}
			continue
		}

		allSectionsForDepartment, err := serverutils.ParseSectionResponse(resp, "")
		if err != nil {
			result <- &resultPayload{
				Dept:         job.Department,
				Error:        err,
				coursesFound: 0,
				execTimeMs:   time.Since(start).Milliseconds(),
			}
			continue
		}

		log.Infof("found %d total CRNs", len(allSectionsForDepartment))

		// Filter down to the minimum set of courses for the department
		filteredCourses := make(map[string]*course)
		for _, section := range allSectionsForDepartment {
			processedSection, exists := filteredCourses[section.CourseNumber]
			if !exists {
				filteredCourses[section.CourseNumber] = &course{
					Number:        section.CourseNumber,
					Name:          section.CourseName,
					TotalSections: 1,
				}

			} else {
				processedSection.TotalSections = processedSection.TotalSections + 1
				filteredCourses[section.CourseNumber] = processedSection
			}
		}

		// Set the data
		collectionRef := clients.fbClient.Collection("departments").Doc(job.Department).Collection(job.Term)
		for _, course := range filteredCourses {
			batch.Set(collectionRef.Doc(course.Number), map[string]interface{}{
				"title":       course.Name,
				"numSections": course.TotalSections,
			}, firestore.MergeAll)
		}
		if len(filteredCourses) != 0 {
			_, err = batch.Commit(ctx)
			if err != nil {
				log.WithError(err).Error("Failed to commit batch")
				result <- &resultPayload{
					Dept:         job.Department,
					Error:        nil,
					coursesFound: len(filteredCourses),
					execTimeMs:   time.Since(start).Milliseconds(),
				}
			}
		}

		result <- &resultPayload{
			Dept:         job.Department,
			Error:        nil,
			coursesFound: len(filteredCourses),
			execTimeMs:   time.Since(start).Milliseconds(),
		}
	}
}

func fetchAndUpdateDepartments(ctx context.Context, batch *firestore.WriteBatch, clients *clients) ([]*models.Department, error) {
	response, err := serverutils.MakeAtlasDepartmentRequest(clients.HttpClient)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Request to myInfo failed")
		return nil, err
	}
	defer response.Body.Close()

	start := time.Now()
	departments, err := serverutils.ParseDepartmentResponse(response)
	elapsed := time.Since(start)
	log.WithContext(ctx).Infof("Scrape time: %v", elapsed.String())
	if err != nil {
		log.WithError(err).WithContext(ctx).Errorf("Department Scrape Failed")
		return nil, err
	}

	departmentsRef := clients.fbClient.Collection("departments")

	for _, dept := range departments {
		deptRef := departmentsRef.Doc(dept.Id)
		batch.Set(deptRef, map[string]interface{}{
			"name":        dept.Name,
			"updatedTime": firestore.ServerTimestamp,
		}, firestore.MergeAll)
	}

	return departments, nil

}
