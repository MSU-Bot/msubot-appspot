package dstore

import (
	"context"
	"errors"

	"cloud.google.com/go/firestore"
	"github.com/SpencerCornish/msubot-appspot/server/models"
	log "github.com/sirupsen/logrus"
)

type fbDStore struct {
	fbClient firestore.Client
}

// New creates a new Firebase implementation of dStore,
func New(fbClient firestore.Client) DStore {
	return fbDStore{fbClient: fbClient}
}

func (f fbDStore) GetMeta(ctx context.Context) (*models.Meta, error) {
	data, err := f.fbClient.Doc("global/global").Get(ctx)
	if err != nil {
		return nil, err
	}
	metaModel := &models.Meta{}
	err = data.DataTo(metaModel)
	if err != nil {
		return nil, err
	}

	return metaModel, nil
}

func (f fbDStore) GetDepartments(ctx context.Context) ([]string, error) {
	emptyDocs, err := f.fbClient.Collection("departments").Select().Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	departmentIDs := make([]string, len(emptyDocs))
	for i := 0; i < len(emptyDocs); i++ {
		departmentIDs[i] = emptyDocs[i].Ref.ID
	}

	return departmentIDs, nil
}

func (f fbDStore) GetCoursesForDepartment(ctx context.Context, term, department string) ([]models.DepartmentCourses, error) {
	data, err := f.fbClient.Collection("departments").Doc(department).Collection(term).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	courses := make([]models.DepartmentCourses, len(data))
	for i := 0; i < len(data); i++ {
		courseData := data[i].Data()
		courses[i] = models.DepartmentCourses{CourseID: data[i].Ref.ID, Title: courseData["title"].(string)}
	}

	return courses, nil
}
func (f fbDStore) GetTrackedSectionByID(ctx context.Context, ID string) (*models.TrackedSectionRecord, error) {
	data, err := f.fbClient.Collection("sections_tracked").Doc(ID).Get(ctx)
	if err != nil {
		return nil, err
	}

	trackedSection := &models.TrackedSectionRecord{}
	err = data.DataTo(trackedSection)
	trackedSection.ID = data.Ref.ID
	if err != nil {
		return nil, err
	}

	return trackedSection, nil
}
func (f fbDStore) GetTrackedSection(ctx context.Context, term, departmentAbbr, courseNumber string) (*models.TrackedSectionRecord, error) {
	data, err := f.fbClient.
		Collection("sections_tracked").
		Where("courseNumber", "==", courseNumber).
		Where("departmentAbbr", "==", departmentAbbr).
		Where("term", "==", term).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	if len(data) > 1 {
		log.WithContext(ctx).WithFields(log.Fields{
			"term":           term,
			"departmentAbbr": departmentAbbr,
			"courseNumber":   courseNumber,
		}).Error("Found duplicate tracked section!")
	}

	trackedSection := &models.TrackedSectionRecord{}
	err = data[0].DataTo(trackedSection)
	trackedSection.ID = data[0].Ref.ID
	if err != nil {
		return nil, err
	}

	return trackedSection, nil
}

func (f fbDStore) UpdateSection(ctx context.Context, sectionID string, section models.Section) error {
	_, err := f.fbClient.Collection("sections_tracked").Doc(sectionID).Set(ctx, map[string]interface{}{
		"courseName":     section.CourseName,
		"courseNumber":   section.CourseNumber,
		"crn":            section.Crn,
		"department":     section.DeptName,
		"departmentAbbr": section.DeptAbbr,
		"instructor":     section.Instructor,
		"openSeats":      section.AvailableSeats,
		"sectionNumber":  section.SectionNumber,
		"term":           section.Term,
		"totalSeats":     section.TotalSeats,
	}, firestore.MergeAll)

	return err
}

func (f fbDStore) UpdateTrackedSection(ctx context.Context, trackedSection models.TrackedSectionRecord) error {
	_, err := f.fbClient.Collection("sections_tracked").Doc(trackedSection.ID).Set(ctx, trackedSection, firestore.MergeAll)
	return err
}

func (f fbDStore) GetTrackedSectionsForUser(ctx context.Context, uid string) ([]models.TrackedSectionRecord, error) {
	data, err := f.fbClient.Collection("sections_tracked").Where("users", "array-contains", uid).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	return trackedSectionDocsToModels(data)
}

func (f fbDStore) GetAllTrackedSections(ctx context.Context) ([]models.TrackedSectionRecord, error) {
	data, err := f.fbClient.Collection("sections_tracked").Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	return trackedSectionDocsToModels(data)
}

func (f fbDStore) GetTrackedSectionsBeforeTerm(ctx context.Context, termCondition string) ([]models.TrackedSectionRecord, error) {
	data, err := f.fbClient.Collection("sections_tracked").Where("term", "<=", termCondition).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	return trackedSectionDocsToModels(data)
}

func (f fbDStore) GetUser(ctx context.Context, userID string) (*models.UserRecord, error) {
	doc, err := f.fbClient.Collection("users").Doc(userID).Get(ctx)
	if err != nil {
		return nil, err
	}
	user, err := userDocToModel(doc)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (f fbDStore) MoveTrackedSectionsToArchive(ctx context.Context, sectionIDs []string) error {
	writeBatch := f.fbClient.Batch()
	writesMade := false

	for i := 0; i < len(sectionIDs); i++ {
		section, err := f.fbClient.
			Collection("sections_tracked").
			Doc(sectionIDs[i]).
			Get(ctx)

		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Failed to get tracked section to move to archive")
			continue
		}

		trackedSections, err := trackedSectionDocsToModels([]*firestore.DocumentSnapshot{section})
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Failed to convert tracked section database entities to models")
			continue
		}
		trackedSectionToMove := trackedSections[0]

		existingArchiveRecords, err := f.fbClient.
			Collection("sections_archive").
			Where("courseNumber", "==", trackedSectionToMove.CourseNumber).
			Where("departmentAbbr", "==", trackedSectionToMove.DepartmentAbbr).
			Where("term", "==", trackedSectionToMove.Term).
			Documents(ctx).GetAll()
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Failed to get existing archive documents")
			continue
		}

		// existing archive record, just add users to it
		if len(existingArchiveRecords) > 0 {
			writeBatch.Update(existingArchiveRecords[0].Ref, []firestore.Update{
				{
					Path:  "users",
					Value: firestore.ArrayUnion(trackedSectionToMove.Users),
				},
			})
			writesMade = true
			// We need a new archive record
		} else {
			newDoc := f.fbClient.Collection("sections_archive").NewDoc()
			writeBatch.Set(newDoc, trackedSectionToMove)
			writesMade = true
		}
	}
	// Commit if writes were made
	if writesMade {
		_, err := writeBatch.Commit(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f fbDStore) AddUserToExistingTrackedSection(ctx context.Context, userUID, sectionID string) error {
	_, err := f.fbClient.Collection("sections_tracked").Doc(sectionID).Update(ctx, []firestore.Update{
		{
			Path:  "users",
			Value: firestore.ArrayUnion(userUID),
		},
	})

	return err
}

func (f fbDStore) AddNewTrackedSection(ctx context.Context, sectionRecord models.TrackedSectionRecord) (*models.TrackedSectionRecord, error) {
	existingData, err := f.fbClient.
		Collection("sections_tracked").
		Where("courseNumber", "==", sectionRecord.CourseNumber).
		Where("departmentAbbr", "==", sectionRecord.DepartmentAbbr).
		Where("term", "==", sectionRecord.Term).
		Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	// Check for existing records, don't make a dupe if there already is one!
	if len(existingData) > 0 {
		return nil, errors.New("tracked Section already exists, not adding a duplicate")
	}

	newDocRef := f.fbClient.Collection("sections_tracked").NewDoc()
	_, err = newDocRef.Create(ctx, sectionRecord)
	if err != nil {
		return nil, err
	}
	sectionRecord.ID = newDocRef.ID

	return &sectionRecord, nil
}

func trackedSectionDocsToModels(data []*firestore.DocumentSnapshot) ([]models.TrackedSectionRecord, error) {
	trackedSections := make([]models.TrackedSectionRecord, len(data))
	for i := 0; i < len(data); i++ {
		section := &models.TrackedSectionRecord{}
		err := data[i].DataTo(section)
		if err != nil {
			// TODO: Decide on a better error
			return nil, err
		}
		trackedSections[i] = *section
	}

	return trackedSections, nil
}

func userDocToModel(data *firestore.DocumentSnapshot) (*models.UserRecord, error) {
	user := &models.UserRecord{}
	err := data.DataTo(user)
	if err != nil {
		// TODO: Decide on a better error
		return nil, err
	}

	return user, nil
}
