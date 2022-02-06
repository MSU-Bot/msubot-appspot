// Package api provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package api

import (
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/labstack/echo/v4"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Scrapes and notifies users of open seats
	// (GET /cron/checktrackedsections)
	CheckTrackedSections(ctx echo.Context) error
	// Cleans up any stale tracked sections
	// (GET /cron/prunetrackedsections)
	PruneTrackedSections(ctx echo.Context) error
	// Gets the courses for a department
	// (GET /department/courses)
	GetCoursesForDepartment(ctx echo.Context, params GetCoursesForDepartmentParams) error
	// Gets the list of departments at MSU
	// (GET /departments)
	GetDepartments(ctx echo.Context) error
	// Gets general info about the web app
	// (GET /meta)
	GetMeta(ctx echo.Context) error
	// Gets the sections for a course
	// (GET /sections)
	GetSections(ctx echo.Context, params GetSectionsParams) error
	// Receives SMS data
	// (POST /service/sms/receive)
	ReceiveSMS(ctx echo.Context) error
	// Gets user data for the specified user
	// (GET /users/{userID})
	GetUserData(ctx echo.Context, userID string) error
	// Updates or sets userdata for the user
	// (PUT /users/{userID})
	UpdateUserData(ctx echo.Context, userID string) error
	// Removes the user from the specified section
	// (DELETE /users/{userID}/section/{sectionID})
	RemoveTrackedSectionForUser(ctx echo.Context, userID string, sectionID string) error
	// Gets tracked sections for the specified user
	// (GET /users/{userID}/sections)
	GetTrackedSectionsForUser(ctx echo.Context, userID string) error
	// Adds tracked sections for the specified user
	// (PUT /users/{userID}/sections)
	AddTrackedSectionsForUser(ctx echo.Context, userID string) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// CheckTrackedSections converts echo context to params.
func (w *ServerInterfaceWrapper) CheckTrackedSections(ctx echo.Context) error {
	var err error

	ctx.Set("appengineApiAuth.Scopes", []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CheckTrackedSections(ctx)
	return err
}

// PruneTrackedSections converts echo context to params.
func (w *ServerInterfaceWrapper) PruneTrackedSections(ctx echo.Context) error {
	var err error

	ctx.Set("appengineApiAuth.Scopes", []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PruneTrackedSections(ctx)
	return err
}

// GetCoursesForDepartment converts echo context to params.
func (w *ServerInterfaceWrapper) GetCoursesForDepartment(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetCoursesForDepartmentParams
	// ------------- Required query parameter "term" -------------

	err = runtime.BindQueryParameter("form", true, true, "term", ctx.QueryParams(), &params.Term)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter term: %s", err))
	}

	// ------------- Required query parameter "deptAbbr" -------------

	err = runtime.BindQueryParameter("form", true, true, "deptAbbr", ctx.QueryParams(), &params.DeptAbbr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter deptAbbr: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetCoursesForDepartment(ctx, params)
	return err
}

// GetDepartments converts echo context to params.
func (w *ServerInterfaceWrapper) GetDepartments(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetDepartments(ctx)
	return err
}

// GetMeta converts echo context to params.
func (w *ServerInterfaceWrapper) GetMeta(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetMeta(ctx)
	return err
}

// GetSections converts echo context to params.
func (w *ServerInterfaceWrapper) GetSections(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetSectionsParams
	// ------------- Required query parameter "term" -------------

	err = runtime.BindQueryParameter("form", true, true, "term", ctx.QueryParams(), &params.Term)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter term: %s", err))
	}

	// ------------- Required query parameter "deptAbbr" -------------

	err = runtime.BindQueryParameter("form", true, true, "deptAbbr", ctx.QueryParams(), &params.DeptAbbr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter deptAbbr: %s", err))
	}

	// ------------- Required query parameter "course" -------------

	err = runtime.BindQueryParameter("form", true, true, "course", ctx.QueryParams(), &params.Course)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter course: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetSections(ctx, params)
	return err
}

// ReceiveSMS converts echo context to params.
func (w *ServerInterfaceWrapper) ReceiveSMS(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.ReceiveSMS(ctx)
	return err
}

// GetUserData converts echo context to params.
func (w *ServerInterfaceWrapper) GetUserData(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "userID" -------------
	var userID string

	err = runtime.BindStyledParameter("simple", false, "userID", ctx.Param("userID"), &userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter userID: %s", err))
	}

	ctx.Set("bearerAuth.Scopes", []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetUserData(ctx, userID)
	return err
}

// UpdateUserData converts echo context to params.
func (w *ServerInterfaceWrapper) UpdateUserData(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "userID" -------------
	var userID string

	err = runtime.BindStyledParameter("simple", false, "userID", ctx.Param("userID"), &userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter userID: %s", err))
	}

	ctx.Set("bearerAuth.Scopes", []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.UpdateUserData(ctx, userID)
	return err
}

// RemoveTrackedSectionForUser converts echo context to params.
func (w *ServerInterfaceWrapper) RemoveTrackedSectionForUser(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "userID" -------------
	var userID string

	err = runtime.BindStyledParameter("simple", false, "userID", ctx.Param("userID"), &userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter userID: %s", err))
	}

	// ------------- Path parameter "sectionID" -------------
	var sectionID string

	err = runtime.BindStyledParameter("simple", false, "sectionID", ctx.Param("sectionID"), &sectionID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter sectionID: %s", err))
	}

	ctx.Set("bearerAuth.Scopes", []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.RemoveTrackedSectionForUser(ctx, userID, sectionID)
	return err
}

// GetTrackedSectionsForUser converts echo context to params.
func (w *ServerInterfaceWrapper) GetTrackedSectionsForUser(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "userID" -------------
	var userID string

	err = runtime.BindStyledParameter("simple", false, "userID", ctx.Param("userID"), &userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter userID: %s", err))
	}

	ctx.Set("bearerAuth.Scopes", []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetTrackedSectionsForUser(ctx, userID)
	return err
}

// AddTrackedSectionsForUser converts echo context to params.
func (w *ServerInterfaceWrapper) AddTrackedSectionsForUser(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "userID" -------------
	var userID string

	err = runtime.BindStyledParameter("simple", false, "userID", ctx.Param("userID"), &userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter userID: %s", err))
	}

	ctx.Set("bearerAuth.Scopes", []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.AddTrackedSectionsForUser(ctx, userID)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/cron/checktrackedsections", wrapper.CheckTrackedSections)
	router.GET(baseURL+"/cron/prunetrackedsections", wrapper.PruneTrackedSections)
	router.GET(baseURL+"/department/courses", wrapper.GetCoursesForDepartment)
	router.GET(baseURL+"/departments", wrapper.GetDepartments)
	router.GET(baseURL+"/meta", wrapper.GetMeta)
	router.GET(baseURL+"/sections", wrapper.GetSections)
	router.POST(baseURL+"/service/sms/receive", wrapper.ReceiveSMS)
	router.GET(baseURL+"/users/:userID", wrapper.GetUserData)
	router.PUT(baseURL+"/users/:userID", wrapper.UpdateUserData)
	router.DELETE(baseURL+"/users/:userID/section/:sectionID", wrapper.RemoveTrackedSectionForUser)
	router.GET(baseURL+"/users/:userID/sections", wrapper.GetTrackedSectionsForUser)
	router.PUT(baseURL+"/users/:userID/sections", wrapper.AddTrackedSectionsForUser)

}
