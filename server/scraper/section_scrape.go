package scraper

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/SpencerCornish/msubot-appspot/server/serverutils"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// HandleRequest scrapes
func HandleRequest(ctx echo.Context, term, dept, course string) error {

	rCtx := ctx.Request().Context()
	client := http.DefaultClient

	if len(course) == 0 || len(dept) == 0 || len(term) == 0 {
		log.WithContext(rCtx).Errorf("Malformed request to API")
		return errors.New("bad syntax. Missing params")
	}
	log.WithContext(rCtx).Debugf("term: %s", term)
	log.WithContext(rCtx).Debugf("dept: %s", dept)
	log.WithContext(rCtx).Debugf("course: %s", course)

	response, err := serverutils.MakeAtlasSectionRequest(client, term, dept, course)

	if err != nil {
		log.WithContext(rCtx).WithError(err).Error("Request to myInfo failed")
		errorStr := fmt.Sprintf("Request to myInfo failed with error: %v", err)
		return errors.New(errorStr)
	}

	start := time.Now()
	sections, err := serverutils.ParseSectionResponse(response, term, "")
	elapsed := time.Since(start)
	log.WithContext(rCtx).WithField("time", elapsed.String()).Info("Scrape Complete")

	if err != nil {
		log.WithError(err).WithContext(rCtx).Errorf("Course Scrape Failed")
		return err
	}

	//Set response headers
	//FIXME: Bruh don't be a jerk
	ctx.Response().Writer.Header().Add("Access-Control-Allow-Origin", "*")
	ctx.Response().Writer.Header().Add("Access-Control-Allow-Methods", "GET")
	ctx.Response().Writer.Header().Add("Content-Type", "application/json")

	return ctx.JSON(http.StatusOK, sections)
}
