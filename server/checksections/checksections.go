package checksections

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/SpencerCornish/msubot-appspot/server/dstore"
	"github.com/SpencerCornish/msubot-appspot/server/models"
	"github.com/SpencerCornish/msubot-appspot/server/serverutils"
	"github.com/labstack/echo/v4"

	log "github.com/sirupsen/logrus"
)

// HandleRequest runs often to check for open seats
func HandleRequest(ctx echo.Context, ds dstore.DStore) error {
	writer := ctx.Response().Writer
	client := http.DefaultClient

	rCtx := ctx.Request().Context()
	defer rCtx.Done()

	// Get the list of sections we are actively tracking
	trackedSections, err := ds.GetAllTrackedSections(rCtx)
	if err != nil {
		log.WithContext(rCtx).WithError(err).Error("Could not get all tracked sections")
		writer.WriteHeader(500)
		return errors.New("could not get all tracked sections")
	}

	// This is the number of concurrent URLFetches that we will do.
	numWorkers := 10

	// A queue of all sections to check
	jobQueue := make(chan *models.TrackedSectionRecord, len(trackedSections))

	// A return channel to let us know a job has completed
	requestCompleteChannel := make(chan int, len(trackedSections))

	// Start up some workers
	for r := 0; r < numWorkers; r++ {
		go sectionCheckWorker(rCtx, jobQueue, requestCompleteChannel, client, ds)
	}

	// Add all sections to the queue
	for _, doc := range trackedSections {
		jobQueue <- &doc
	}

	close(jobQueue)

	// Wait for the jobs to finish
	for i := 0; i < len(trackedSections); i++ {
		<-requestCompleteChannel
	}

	writer.WriteHeader(200)
	return nil
}

func sectionCheckWorker(ctx context.Context, jobs <-chan *models.TrackedSectionRecord, returnChannel chan<- int, client *http.Client, ds dstore.DStore) {
	for record := range jobs {
		// Make a request to Atlas
		resp, err := serverutils.MakeAtlasSectionRequest(client, record.Term, record.DepartmentAbbr, record.CourseNumber)
		if err != nil {
			log.WithContext(ctx).Errorf("Making Atlas request failed for record ID %s", record.ID, err)

			returnChannel <- 0
			continue
		}

		// Parse into a section struct
		newSectionData, err := serverutils.ParseSectionResponse(resp, record.Term, record.Crn)
		if err != nil {
			log.WithContext(ctx).Errorf("Parsing section failed: %v", err)

			returnChannel <- 0
			continue
		}
		// If we somehow get back more than one section, something super borked earlier
		if len(newSectionData) > 1 {
			log.WithContext(ctx).Errorf("Something went wrong with parsing the section response. Expected 1 section, recieved %v", len(newSectionData))

			returnChannel <- 0
			continue
		}

		// If we didn't get back any, warn us and move on.
		// This typically occurs when Banner is down
		if len(newSectionData) == 0 {
			log.WithContext(ctx).Warningf("Couldn't find section from MSU ID: %s", record.ID)

			returnChannel <- 0
			continue
		}

		// Parse the new available seats to an int
		newSeatsAvailable, err := strconv.Atoi(newSectionData[0].AvailableSeats)
		if err != nil {
			log.WithContext(ctx).WithError(err).Errorf("couldn't parse newSeatsAvailable")

			returnChannel <- 0
			continue
		}

		if len(record.Users) < 1 {
			log.WithContext(ctx).Infof("Record %s has %d users. Deleting CRN", record.ID, len(record.Users))

			// TODO: Bulk operation for greater efficiency
			err := ds.MoveTrackedSectionsToArchive(ctx, []string{record.ID})
			if err != nil {
				log.WithContext(ctx).Errorf("Failed to move the stale section data: %v", err)
			}
			returnChannel <- 0
			continue
		}

		// If there are seats available
		if newSeatsAvailable > 0 {
			log.WithContext(ctx).Infof("The resource %s has %d open seats. Sending a message to %d users.", record.ID, newSeatsAvailable, len(record.Users))

			sendOpenSeatMessages(ctx, client, ds, record.Users, newSectionData[0])

			//TODO: Bulk this
			err := ds.MoveTrackedSectionsToArchive(ctx, []string{record.ID})
			if err != nil {
				log.WithContext(ctx).Errorf("Failed to move the stale section data: %v", err)
			}

			returnChannel <- 0
			continue
		}

		// If we get here, we just need to update the stored section model so it's all clean and nice
		err = ds.UpdateSection(ctx, record.ID, newSectionData[0])
		if err != nil {
			log.WithContext(ctx).WithError(err).Errorf("couldn't update section data")

			returnChannel <- 0
			continue
		}
		returnChannel <- 0
	}
}

// TODO: Move to messager package
func sendOpenSeatMessages(ctx context.Context, client *http.Client, ds dstore.DStore, users []string, section models.Section) error {
	var userNumbers string
	message := fmt.Sprintf("%v%v - %v with CRN %v has %v open seats! Get to MyInfo and register before it's gone!", section.DeptAbbr, section.CourseNumber, section.CourseName, section.Crn, section.AvailableSeats)
	for _, user := range users {
		userRecord, err := ds.GetUser(ctx, user)
		if err != nil {
			log.WithContext(ctx).Errorf("Unable to get user %s", user)
		}

		if userNumbers == "" {
			userNumbers = userRecord.PhoneNumber
		} else {
			userNumbers = fmt.Sprintf("%v<%v", userNumbers, userRecord.PhoneNumber)
		}
	}
	resp, err := serverutils.SendText(client, userNumbers, message)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error sending texts")
		return err
	}
	log.WithContext(ctx).WithField("response", resp).Info("Texts sent")
	return nil
}
