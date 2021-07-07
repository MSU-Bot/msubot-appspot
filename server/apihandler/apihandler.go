package apihandler

import (
	"net/http"

	"github.com/SpencerCornish/msubot-appspot/server/datastore"
	"github.com/SpencerCornish/msubot-appspot/server/dstore"
)

type ApiHandler interface {
	GetMeta(w http.ResponseWriter, r *http.Request)
	GetDepartments(w http.ResponseWriter, r *http.Request)
	GetCoursesForDept(w http.ResponseWriter, r *http.Request)
	GetSections(w http.ResponseWriter, r *http.Request)

	CheckTrackedSections(w http.ResponseWriter, r *http.Request)
	PruneTrackedSections(w http.ResponseWriter, r *http.Request)

	ReceiveSMS(w http.ResponseWriter, r *http.Request)
	HealthCheck(w http.ResponseWriter, r *http.Request)

	GetUserSections(w http.ResponseWriter, r *http.Request)
	TrackSections(w http.ResponseWriter, r *http.Request)
	RemoveTrackedSections(w http.ResponseWriter, r *http.Request)
}

type apiHandler struct {
	datastore dstore.DStore
}

func New(ds datastore.DataStore) ApiHandler {
	return apiHandler{datastore: ds}
}

func (a apiHandler) GetMeta(w http.ResponseWriter, r *http.Request) {

}

func (a apiHandler) GetDepartments(w http.ResponseWriter, r *http.Request) {

}
func (a apiHandler) GetCoursesForDept(w http.ResponseWriter, r *http.Request) {

}
func (a apiHandler) GetSections(w http.ResponseWriter, r *http.Request) {

}

func (a apiHandler) CheckTrackedSections(w http.ResponseWriter, r *http.Request) {

}
func (a apiHandler) PruneTrackedSections(w http.ResponseWriter, r *http.Request) {

}

func (a apiHandler) ReceiveSMS(w http.ResponseWriter, r *http.Request) {

}
func (a apiHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {

}

func (a apiHandler) GetUserSections(w http.ResponseWriter, r *http.Request) {

}
func (a apiHandler) TrackSections(w http.ResponseWriter, r *http.Request) {

}
func (a apiHandler) RemoveTrackedSections(w http.ResponseWriter, r *http.Request) {

}
