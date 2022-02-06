package pruner

import (
	"fmt"
	"time"

	"github.com/SpencerCornish/msubot-appspot/server/dstore"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// HandleRequest is run daily to clean up expired course checkers from old semesters.
func HandleRequest(ctx echo.Context, ds dstore.DStore) error {
	writer := ctx.Response().Writer

	rCtx := ctx.Request().Context()
	defer rCtx.Done()

	term := getPruneTerm()

	log.WithContext(rCtx).Infof("Removing all trackedsections where term <= %s", term)
	// Get the list of sections we are actively tracking
	oldTrackedSections, err := ds.GetTrackedSectionsBeforeTerm(rCtx, term)
	if err != nil {
		log.WithContext(rCtx).WithError(err).Error("Unable to retrieve old tracked sections")
		writer.WriteHeader(500)
		return err
	}

	log.WithContext(rCtx).Infof("Number of expired courses: %d", len(oldTrackedSections))

	ids := make([]string, len(oldTrackedSections))
	for i, section := range oldTrackedSections {
		ids[i] = section.ID
	}

	return ds.MoveTrackedSectionsToArchive(rCtx, ids)
}

func getPruneTerm() string {
	now := time.Now()
	term := "00"
	year := now.Year()
	if now.Month() > 10 {
		// If our current month is after October or greater, remove fall (year) and before
		term = fmt.Sprintf("%d%d", year, 70)

	} else if now.Month() > 8 {
		// If our current month is September or greater, we should remove summer (year) and before
		term = fmt.Sprintf("%d%d", year, 50)

	} else if now.Month() > 3 {
		// If our current month is April or greater, remove spring (year) and before
		term = fmt.Sprintf("%d%d", year, 30)
	}
	return term
}
