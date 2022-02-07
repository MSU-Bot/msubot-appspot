package tracksections

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"regexp"

	"cloud.google.com/go/firestore"
	"github.com/SpencerCornish/msubot-appspot/server/mauth"
	"github.com/SpencerCornish/msubot-appspot/server/models"
	"github.com/SpencerCornish/msubot-appspot/server/serverutils"
	log "github.com/sirupsen/logrus"
)

// HandleRequest scrapes
func HandleRequest(w http.ResponseWriter, r *http.Request) {

}

func getSectionMetadata(ctx context.Context, request trackRequest, crn string) (models.Section, error) {
	client := http.DefaultClient

	resp, err := serverutils.MakeAtlasSectionRequest(client, request.Term, request.DepartmentAbbr, request.Course)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error getting new section data from ATLAS")
		return models.Section{}, err
	}
	sectionDatas, err := serverutils.ParseSectionResponse(resp, request.Term, crn)
	if err != nil {
		return models.Section{}, err
	}

	return sectionDatas[0], nil

}

// decodeRequest decodes and validates the request
func decodeRequest(body io.ReadCloser, request *trackRequest) error {

	requestDecoder := json.NewDecoder(body)
	requestDecoder.DisallowUnknownFields()

	err := requestDecoder.Decode(&request)
	if err != nil {
		return err
	}

	if !isValidTerm || err != nil || len(request.Crns) == 0 {
		return err
	}

	return nil
}
