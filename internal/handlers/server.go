// Package handlers implements the OpenAPI strict server. Use logging.From(ctx) for request-scoped slog.
package handlers

import (
	"context"
	"net/http"

	"github.com/MSU-Bot/msubot-appspot/internal/api"
)

// Server implements api.StrictServerInterface. Replace stubs with real behavior.
type Server struct{}

func stubErr() api.ErrorResponse {
	return api.ErrorResponse{Code: "NOT_IMPLEMENTED", Message: "handler not implemented yet"}
}

func (Server) AdminScrapeSections(_ context.Context, _ api.AdminScrapeSectionsRequestObject) (api.AdminScrapeSectionsResponseObject, error) {
	return api.AdminScrapeSections202Response{}, nil
}

func (Server) AdminUpdateIndex(_ context.Context, _ api.AdminUpdateIndexRequestObject) (api.AdminUpdateIndexResponseObject, error) {
	return api.AdminUpdateIndex202Response{}, nil
}

func (Server) Login(_ context.Context, _ api.LoginRequestObject) (api.LoginResponseObject, error) {
	return api.LogindefaultJSONResponse{Body: stubErr(), StatusCode: http.StatusNotImplemented}, nil
}

func (Server) Register(_ context.Context, _ api.RegisterRequestObject) (api.RegisterResponseObject, error) {
	return api.RegisterdefaultJSONResponse{Body: stubErr(), StatusCode: http.StatusNotImplemented}, nil
}

func (Server) DeleteSession(_ context.Context, _ api.DeleteSessionRequestObject) (api.DeleteSessionResponseObject, error) {
	return api.DeleteSessiondefaultJSONResponse{Body: stubErr(), StatusCode: http.StatusNotImplemented}, nil
}

func (Server) ListCourses(_ context.Context, _ api.ListCoursesRequestObject) (api.ListCoursesResponseObject, error) {
	return api.ListCoursesdefaultJSONResponse{Body: stubErr(), StatusCode: http.StatusNotImplemented}, nil
}

func (Server) PostCoursesOccupancy(_ context.Context, _ api.PostCoursesOccupancyRequestObject) (api.PostCoursesOccupancyResponseObject, error) {
	return api.PostCoursesOccupancydefaultJSONResponse{Body: stubErr(), StatusCode: http.StatusNotImplemented}, nil
}

func (Server) ListCourseSuggestions(_ context.Context, _ api.ListCourseSuggestionsRequestObject) (api.ListCourseSuggestionsResponseObject, error) {
	return api.ListCourseSuggestionsdefaultJSONResponse{Body: stubErr(), StatusCode: http.StatusNotImplemented}, nil
}

func (Server) GetCourseByCrn(_ context.Context, _ api.GetCourseByCrnRequestObject) (api.GetCourseByCrnResponseObject, error) {
	return api.GetCourseByCrndefaultJSONResponse{Body: stubErr(), StatusCode: http.StatusNotImplemented}, nil
}

func (Server) GetStats(_ context.Context, _ api.GetStatsRequestObject) (api.GetStatsResponseObject, error) {
	return api.GetStatsdefaultJSONResponse{Body: stubErr(), StatusCode: http.StatusNotImplemented}, nil
}

func (Server) DeleteMe(_ context.Context, _ api.DeleteMeRequestObject) (api.DeleteMeResponseObject, error) {
	return api.DeleteMedefaultJSONResponse{Body: stubErr(), StatusCode: http.StatusNotImplemented}, nil
}

func (Server) GetMe(_ context.Context, _ api.GetMeRequestObject) (api.GetMeResponseObject, error) {
	return api.GetMedefaultJSONResponse{Body: stubErr(), StatusCode: http.StatusNotImplemented}, nil
}

func (Server) PatchMe(_ context.Context, _ api.PatchMeRequestObject) (api.PatchMeResponseObject, error) {
	return api.PatchMedefaultJSONResponse{Body: stubErr(), StatusCode: http.StatusNotImplemented}, nil
}

func (Server) ListMySections(_ context.Context, _ api.ListMySectionsRequestObject) (api.ListMySectionsResponseObject, error) {
	return api.ListMySectionsdefaultJSONResponse{Body: stubErr(), StatusCode: http.StatusNotImplemented}, nil
}

func (Server) AddMySection(_ context.Context, _ api.AddMySectionRequestObject) (api.AddMySectionResponseObject, error) {
	return api.AddMySectiondefaultJSONResponse{Body: stubErr(), StatusCode: http.StatusNotImplemented}, nil
}

func (Server) DeleteMySection(_ context.Context, _ api.DeleteMySectionRequestObject) (api.DeleteMySectionResponseObject, error) {
	return api.DeleteMySectiondefaultJSONResponse{Body: stubErr(), StatusCode: http.StatusNotImplemented}, nil
}
