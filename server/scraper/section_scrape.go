package scraper

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/SpencerCornish/msubot-appspot/server/models"
	"github.com/SpencerCornish/msubot-appspot/server/serverutils"
	log "github.com/sirupsen/logrus"
)

// HandleRequest scrapes
func HandleRequest(ctx context.Context, term, dept, course string) ([]models.Section, error) {

	client := http.DefaultClient

	response, err := serverutils.MakeAtlasSectionRequest(client, term, dept, course)

	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Request to myInfo failed")
		errorStr := fmt.Sprintf("Request to myInfo failed with error: %v", err)
		return nil, errors.New(errorStr)
	}

	start := time.Now()
	sections, err := serverutils.ParseSectionResponse(response, term, "")
	elapsed := time.Since(start)
	log.WithContext(ctx).WithField("time", elapsed.String()).Info("Scrape Complete")
	if err != nil {
		log.WithError(err).WithContext(ctx).Errorf("Course Scrape Failed")
		return nil, err
	}

	return sections, nil
}
