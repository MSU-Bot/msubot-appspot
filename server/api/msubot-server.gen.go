// Package api provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
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

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xaa3fbxtH+K/PizTlNW15ASrJP9am0fGNsXY4otUltnWYJDImNgF14d0GKUfjfe/YC",
	"kCAWMmU7btL2ky0SOzP7zDNX8D6IeJZzhkzJ4Pg+EChzziSaP64ZKVTCBf0Z4xdCcKE/jFFGguaKchYc",
	"B6MoQilB8VtkQCVkVErK5sAFULYgKY2D9boTyCjBjBihp6iI/jcXPEehqFUV8UJIlFeCRLcYN/VcJQiK",
	"K5ICK7IpCuAzcGdA2UMwXcHp5PoZV1BIFDLoBHhHsjzF4PjwKOwEapVjcBxQpnCOIlh3goyrNl14p0Bx",
	"mKIWFgNnoBKEBElslfOZ/YBnmJM5bisLJowvM5QKRWUjEQiML0EWec6Fwvj/gsogqQRlc22P1ionyNR+",
	"AJjHQSIzpj7y7va5vfQ8SvK6Ewj8UFCh3fhu17Gl3u273lRC+PQnjJS27iKlCz45nTSZ8lLwzG+2lvwH",
	"CXnCGTrra2758+AwfHJ0dHRweDj0oX+KUpI5jpnyemCCBmaeK14ooNb/lEU804TP7GHIiIoSlLCkKgFt",
	"iCOKVITFRMRwnqvueaHgFldLLmLZg41gyj5P7pg1xSaY5p8u9bU+vZE5lpDiTME0JewWKAOSpsBVoplO",
	"JMree1aD3IL1ANjX1+PnLd5k9EOBQGNkis4oCphxYYxzxtcUhS9eYPjkzXWWJjP+t0ux/KCWdxdvnhJ2",
	"5/X1BV+iyEl0226A/qbEo3ociJQ8okRhbEHT317xdtI91rArvPNw74QbUpbm+BB4O55ceQXyprgzF9oM",
	"lgmNkm2ZsCQSBEZIFyZc9w+fK501RhkvvOnLpJQoIWKOxpFWheaiVj45ncC35plLohD+BGe8py97zaiS",
	"f6yZ8U3YC8MnrQbo835/OuUkz1MakWmKkKPYhIQ2oWBUNZX5b2s+aOhZ5fiQk2QmfcLMNf1Wb7JwToSS",
	"OuY2PmvEc5vzht483ci6E4ys6t2kSxaEphqzCRJraSW66y8vNu+fkQxrTwcvqUxgooqYooTx2IeGO2pj",
	"qXZ4GA4vLtuPlD7ZHHiWIosNFM0jAmO6c5fgoBdCF85ZShn6D7HagWF4EB4MfPePMVej6XTnAifnp9/7",
	"5FJPG+KqrpYpGEmbmZDq2m8c1tEJnirAOyotR0wSX0mFWZ1/PMN/FtQLCGVSiSJSfMfmC3KntIrveMJ8",
	"51IekZI1m1OX589eDwYD3wFWZNdl81EdOPKB6K7no8Lb0TPoQuhVoMgtsiZTD30qFIrMV+ur/i3GEs8Z",
	"FxlR8OP9D0jE+v4g/OUo/OVp+Mvk7Hz9Yw3mYaiJ4TWN7gbE1dXrwUEYdgeHQ/8JndWalzn4eOtlrrbF",
	"xJ3I8rVd2i/N6GcV/l++6dLnW1sAK7vRCXxWgV1iGvEM/V32MkHTzHABjCvjdW0EJDavUlxgDAScDMCM",
	"0LSvO9ltk5QosFI85TxFwhrOcffulOA2vWH5XwiqVhM9PrlMnOfI5pThKKejQiX6M6pNt6OJFmhybvB9",
	"d1Q+2j0RfCt0SU7f4EpjMUUiUJRi7F8vDcuD4+C7v+t+wkxu5iLm242URKk8WK9N4phxO8MxRSKDqgEm",
	"OA4y/KvMkfWiyObEXQ+PLsYmmdlcF3SClEbIpAkRd49RTqIEYdjT0VGI1Kk+7veXy2WPmG97XMz77qjs",
	"vx2fvDibvOgOe2EvUVlq404ZtlhFG2iCTrBAIa1Fg17YC/XTPEdGcmoLglGcE5UY/PuR4KwfJRjdurnT",
	"JSnz7Rx9g0MkSI7SXLTIY9M72kgEDV0HCIs13TS3pZ2ywPDNssW01DocTZYdx7qMaPVupJqU6jv14X0Y",
	"hh5TCjOvz4o0XYEgDBSRt/rGR76nx2XpkSgWKADNDsA8PfQ8LWWhJ91Im62bEoZLiFIiJWSoSEwUqZE6",
	"OH7no/O7m/VNJ5BFlhGx2kLPA5Iew3NkIE121Hl/Ls3Mqel+o3VZb+WiYLivt050uEoociBspQehFKsN",
	"Q3m4A7KIEu0kzpxjzUX1JKWbs/cs5iaBmIJsHWyetK0bEXrG5sDT2O0YDJ+8rr7Qtv/bXP14d+0Dn99V",
	"MeoON0Om+m5p0OqjS1SFYNLk51vGl6zas2hXENiIch/M6QIZ6ILogfgVqhN7/CUXz6ujJuoFyVCZXuXd",
	"5/QIbf2BydwfChSrTeJ2ZXtTLGw9sSs0jcVOSVt3GpYlXCjQ4sz14+0rNXpRnwlbLcP+Ztz4WRnxap/i",
	"Ri9tZf8naTvGjTwSx1R/RdKLWvfR7Ih26mSjsjhvwhKFRqBgbglZkfQVKkueNtpsMRQXKFacYYOl+9Ez",
	"pdKM7VvnyorXsatFxSHneZHqwXdzgtXWi37WPt+y5TPBpwqzB9EmQpCVD+wRmDUvlOrboPYBQZTGoR3s",
	"zO2KH0S5LC6bpOxymFREUaloJDvV7qgEdLPW9ENrttRfCtNvBM6C4+D/+5tde99txPtG0ReCeY4MhZ4U",
	"2YwDmfLCdrBLnALJ83aY9+5gRldvRxODZXmkAtZm2K22RudBypkf3q1K9r8c+zlm2FQH791K730A345S",
	"yYGyKC1ilHDCBcIVmcsOpPQW4fKsvk4bhIPxZYuJ1ptftQjsFTPlmmqPsCmJZmuB3O6JXGHoBMPw0LMi",
	"5RuOP3DW2wnbzIQxLBNkMFIpkUD1HCuQRIlZPZoY0iOkQKJBaEmatTAjUDmkNZDFgkbYl5nsuzWgmea5",
	"9OZP7SdTf2pLULdLlJ7IvbQyJ6cTRwqU6hmPV1/exdX7H1+9X+/T/56/+ZRmt3KBu6o0kJjZZYO6g9mB",
	"bmaR/r0d6NetObTyqdtqlGULZoJnztlGrD9jXksUz60ZD2ZM/6sTt5LWqmvRPzo4fHLy3Xjyw8touXoq",
	"/nF2MXwVpT8f3d09/UuZFPTou8kJ1d7iN5YTzOZqn4TQCGODip1OO8FhOGjTVN2h33wv3hiUtjcruyOS",
	"YYJZK1n/uwoqc4y0t+LSTSXd7BvTm3UnyAsPsa7NSuGTuWWP/97p9SumojZmrX+ztC6XTF+b2CUTTX/o",
	"OF6jeAuxm1m0bEn79+4/LrPGmKLv/d4lZnzhYkBfwV1+i/9VcDmB3uqmhdR3LS+5uLZG/97CovMoGx0q",
	"NTO73f0NrNz0JSrDTn7TjhTGN7FzaH2fZHu4wUdWYJFAExQMl1D3sVnPkTguM29bS2gMyQUuKC9kunJH",
	"FN8R97WibZvzHyH7o2JOfryF2V3otdSwHhjIskIqmLrAZEpnRO+q8xWqnUXn7zf6fnNzULPtaXViGQRf",
	"rRfak077t0SjOLZEzQVfUBOldWI1NXS264aHnKM4/s8i56/YMT3Ay71Gt930nAok8cq+TNGeqo3iu4n7",
	"v64OWLJ/cghZXXoOtuzdvGKVx/3+ggplv+3JJZnPUSTFtBfxrG+3p/3dF6r9gXtjuuMCe/h1MTXvfEeF",
	"4nDKo1u7ZN7VmsmiO+WqV0Q90SN5LnOujFK/cPdLmQvBYxgzqQiLMNAgubveb7oUO7h7tkSFMOjpm3YV",
	"77rNAGeaMSRNpbGy3IvZ1+k7mzj9lD1Q/VpKk6JCBsyxLTHlL2F39hZl2pF1ObXq6X6LuyWs2gT57dqI",
	"WlKV8EJtyzNkvVn/KwAA///wszCeDC4AAA==",
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file.
func GetSwagger() (*openapi3.Swagger, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromData(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error loading Swagger: %s", err)
	}
	return swagger, nil
}