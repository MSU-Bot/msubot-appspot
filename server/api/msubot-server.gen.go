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
	// (GET /user)
	GetUserData(ctx echo.Context) error
	// Updates or sets userdata for the user
	// (PUT /user)
	UpdateUserData(ctx echo.Context) error
	// Removes the user from the specified section
	// (DELETE /user/section/{sectionID})
	RemoveTrackedSectionForUser(ctx echo.Context, sectionID string) error
	// Gets tracked sections for the specified user
	// (GET /user/sections)
	GetTrackedSectionsForUser(ctx echo.Context) error
	// Adds tracked sections for the specified user
	// (PUT /user/sections)
	AddTrackedSectionsForUser(ctx echo.Context) error
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

	ctx.Set("bearerAuth.Scopes", []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetUserData(ctx)
	return err
}

// UpdateUserData converts echo context to params.
func (w *ServerInterfaceWrapper) UpdateUserData(ctx echo.Context) error {
	var err error

	ctx.Set("bearerAuth.Scopes", []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.UpdateUserData(ctx)
	return err
}

// RemoveTrackedSectionForUser converts echo context to params.
func (w *ServerInterfaceWrapper) RemoveTrackedSectionForUser(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "sectionID" -------------
	var sectionID string

	err = runtime.BindStyledParameter("simple", false, "sectionID", ctx.Param("sectionID"), &sectionID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter sectionID: %s", err))
	}

	ctx.Set("bearerAuth.Scopes", []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.RemoveTrackedSectionForUser(ctx, sectionID)
	return err
}

// GetTrackedSectionsForUser converts echo context to params.
func (w *ServerInterfaceWrapper) GetTrackedSectionsForUser(ctx echo.Context) error {
	var err error

	ctx.Set("bearerAuth.Scopes", []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetTrackedSectionsForUser(ctx)
	return err
}

// AddTrackedSectionsForUser converts echo context to params.
func (w *ServerInterfaceWrapper) AddTrackedSectionsForUser(ctx echo.Context) error {
	var err error

	ctx.Set("bearerAuth.Scopes", []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.AddTrackedSectionsForUser(ctx)
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
	router.GET(baseURL+"/user", wrapper.GetUserData)
	router.PUT(baseURL+"/user", wrapper.UpdateUserData)
	router.DELETE(baseURL+"/user/section/:sectionID", wrapper.RemoveTrackedSectionForUser)
	router.GET(baseURL+"/user/sections", wrapper.GetTrackedSectionsForUser)
	router.PUT(baseURL+"/user/sections", wrapper.AddTrackedSectionsForUser)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xaeXcbtxH/Kug27zVtKXJ1xY3+Ci1fjK3jiVKb1NZrwN0hF9EuAOMgxTr87n049uJi",
	"LSZOXOW1f0nk7gwG8/vNgQE/RAkrOKNAlYxOPkQyyaDA9t9nwLFQBVBlPnHBOAhFwD7Ds5mAJcGKMGo+",
	"pyATQbj7GF1ngKo3IEVzJgrE5khlgNJKK6K4gGgQwT0ueA7RSXR6cfZdNIjUmptPUglCF9FmENkXg6vM",
	"dZ5bPV31bc2s4FqBQNOEAE2gu8pmEAl4r4mANDp5296gt+C2EmKzHyFRxrTaSadMCwldVyX2+8mz8Abc",
	"U0R1MQMxIilQReYEBPoyYVRhQiU6vbh6jihT1pY/t7Z1dBxPzkMeU0TlPS5resut3lI5oUowpBj6Fi8x",
	"mjzoqGp75aIhLz0XgomQa9KAkfZlZJ8NIkMdrKKTiFB1eFBbQ6iCBQijvAAp8aJXUfn44Y3YBcvXbzeD",
	"6AwU7gNUXguc3BnRkJMVUzj3qBpnexmknBCardHZ9OYpU0hLELKJwNFxHNwlU31rwb0ygM3AKEsRoxbb",
	"DHDqFvdgZ6wA7jxRwz2lbFWANKFR2ogFIMpWSGrOmVCQ/iHIMLhXcuqTw8MOsK8jacJesZ+7d/feTuv8",
	"LM1BKlfAlus29xoi92VOlmx6Nu0y5YVgRdhso/lPEvGM0TL4W7D8df8o/ur4+Pjw6Ogg5P0zx9IJVUEE",
	"pmDdzLhiWiHi8Cc0YQWhizIiUIFVkoFEK6IyZAzxRJEK0xSLFF1wtXehFbqD9YqJVA5RrZjQT9M7oV21",
	"GeT8l2t9ZaRrnROJcpgrNMsxvUOEIpzniKnMMB1LkMN3tOVy56yPOPvmpi+Ja0rea0CN7D1nwhpXJ596",
	"ofj5c4i/en1T5Nmc/f1KrN6r1f3l6yeY3gexvmQrEBwnd/0GmCelP6rXEZaSJa4CW6eZp9esn3Q/17Br",
	"uA9w75RZUpbmhDzwZjK9DipkXXXnPrQpWmUkyZo60QpLJCABsrThunv4XJusMS6YDqYvm1KSDIsFWCDd",
	"EoaLZvHp2RR9ad+5wgrQX9A5G5rN3lCiZLs8fxEP4/irXgOMfE9X4BbHnOckwbMcEAdRh4QxQVOiuouF",
	"d2u/6Kyz5vAxkGQhQ8rsNns6iyoLm55ImpirMevEcx94B8E83cm6U0jK5nOrNV1ikhufTQE7SyvVe+Hy",
	"4vL+uW8yaw+8IDJDU6VTAhJNJiFveFEXSy3hg/jg8qpfpMSkFniaA02tK7oiAlKytZfocBijPXRBc0Ih",
	"LERbAgfxYXy4H9p/ClyNZ7OtDfQ14yTQhviqa3QKivNuJiSm9lvABibBE4XgnkjHEZvE11JB0eYfK+Bf",
	"mgQdQqhUQieKbdl8ie+VWeJbltGQXM6S6shSS11dPH21v78fPHno4qZsPiqB45AT/fZCVHgzfor2UBxc",
	"QOE7oF2mHoWWUCCKUK2v+rcUSn+6thn98OF7wGLz4TD+6Tj+6Un80/T8YvNDy80HsSFG+AyxHRDX16/2",
	"D+N4b//oICxhslp3M4cPt152aw0mbkVWqO0yuHSjHwpM8rbV/r9v/N9hkvQg7YH79bs1I9/bOzjdnRbi",
	"kyrzCvKEFRBuz1cZ2C6ICXOqtHQxRqDMJWQCS0gRRl4Hsh4dmRa4aZISGqqFZ4zlgGkHVb/vyrldGF3g",
	"aEHUeppkhsk2hXMOdEEojDkZa5WZ74gx3Z1pyvP4SfTd3rh8de9UsEbMY05ew9r4YgZYgCjVuE8vylPl",
	"t/8wjYgdetiN2Ke1lkwpHm02NuPMmTv8UYUT1aBaVEg9Y+obyYEabtmEuo3y+HJiM6FLlNEgykkC1I0L",
	"/F7GHCcZoIOhCS0tcr/8yWi0Wq2G2D4dMrEYeVE5ejM5fX4+fb53MIyHmSryxsHfZ+TaPdEgWoKQzqL9",
	"YTyMzduMA8WcuGpiF+ZYZRaDUSIYHSUZJHf+0OoznH26gNCpIxGYg7Qb1Ty1jaefbxj3DRCmqaGc4bd0",
	"RzRkOecYY/txE8s2RU9SU4PM8v48Ni2XNxyTnBkXGBsO4jhgik4SkHKu83yNBKZIYXnnoJljnasSSh8g",
	"vtEywqMfpasPbhJm/vtCwDw6if44qkdlIz8nG7nJhqVI2wLwD2qCRydvQ9R+e7u5HURSFwUW64YXA84y",
	"Z3kOFEmbYk3xWEh7cDXUvzVrOdS40BR2Re3UhK5EmiNM1+Y0lUM1piiFB0jqJDNgMeoBTnIspTmOmQ7v",
	"HU2ZTSa2qjug7Zuu/8PCHNQZYnnqBxWWV0HIL43t/2OQ7wJBGO563Dny04tenK9AaUGlzfd3lK1oNfAx",
	"cOLmYNZ9sSBLoMhU5gBML8EPPOULJp41h64cC1yAsk3T209pVvoaFVsJ3msQ67oQ+P6hLj6uPtWAdmZ/",
	"HcsyJtxI2m6/b47smuKQCY3eZXczbsPM3pmnREEhHyJsZ0Zdn6ewEHgd4rIHF61AGIdomhqpo/gocEJn",
	"JZNGLQ55mf9G+FWR9RKUbEy6u1xvhBUsQawZhU5o7RZTOZF26NGQK0v+wA1mFUOccZ1jBQ0J2hrOhkPt",
	"WcOWz8uYXbgyRkuckxSVdj0e1EOYYGUg6ce98EP/jwJuXkqxwnVhlCCWIEzyVkQqkshBNQQssa3n02GU",
	"7XXD54DXLtQdqvxekF0ABYFz21QiPGPaHWJWMEOY835kd25gx9dvxlMLXylSYemKYqOrNaWLMBpGtNHA",
	"/L8sfooZrhyhd34c/C5CX45zyRChSa5TkOiUCUDXeCEHKCd3gK7O26PY/Xh/ctVjYnUL+sjqdjni3CEF",
	"l0Rz9Vo2W+FHVohbEYXrG+jemBVLksBIFnLkp8V26MNkMDsbSGyhbc3K/chZBoL0yumcnk09/iDVU5au",
	"f300q2vCUOLd7HLCuXj9WWHUFO45JOYUH0LUe05aD5ta2ADRo+Yx1H5SF8y5FTH8IKysrGguWOEZY5WF",
	"M+yNBPEMf6a6aUeOu0RjJ/7s7qyTHsGBtDkR2z6KWjjsONCB4Mue5JCQOQG3kQbS7or8djOIuA6ge2PH",
	"QL8YYCfewvg3jNA+eDePllvllO33wq6SDraz8kRr8ayHXWUWKVu40Qf/z+TZxpEuh9Bd6hUUbOnph7XK",
	"vKsa1Kt47RUGS4RR0h5JvWDixpn60b4u/OOA8ocLvrw325S9ve9fJKv1E/HP88uDl0n+7+P7+ydfl30L",
	"xyqr25bKBb9G57IVtsZJwu479c5qj6MeO9GawD+A+A50kw9Xr+1xXU/mHCLr2kJLhWaek1QZfwWHoS9B",
	"bY1Ca+I9rv6zW/F6PaJlee38yMvgjpjuXg3HaerYwgVbkhQqJ017Vxg081aAIeM0/QhDfsNi+RFy7NTM",
	"tq1GOBeA07W7QDDbbp1hcJqWzt4MooN4/4HxfyLA1kUKK7S9EE076o56MiAXsCRMy8oCxbbUPXYeO8b9",
	"Yh67tcSyrG71taQ8GY2WRCj3dChXeLEAkenZMGHFyE25RtuXkKN9f8u4BZ0TfqVn9p50rBVDZyy5c5fZ",
	"26sWUu/NmBrqZCiGmHPJmbKLhpX7n6ZcCpaiCZUK0wQaWnOW4DxjUp38Lf66K/3GPEYKpLLG3FYu+lBX",
	"YXfUCZzItbBONw7aU2zPTwgZNQTFeS6tGeUMwt1eb009zFtOoPpVk+Fg5VBkxRpqyl+sbh0cy5Qh23pa",
	"5cf/ZrahrDqKh+2qVa2IyphWTX02Nm43/wkAAP//ZMFBQ18wAAA=",
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
