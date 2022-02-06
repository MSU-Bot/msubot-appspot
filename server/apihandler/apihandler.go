package apihandler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/SpencerCornish/msubot-appspot/server/checksections"
	"github.com/SpencerCornish/msubot-appspot/server/dstore"
	"github.com/SpencerCornish/msubot-appspot/server/gen/api"
	"github.com/SpencerCornish/msubot-appspot/server/mauth"
	"github.com/SpencerCornish/msubot-appspot/server/messenger"
	"github.com/SpencerCornish/msubot-appspot/server/pruner"
	"github.com/SpencerCornish/msubot-appspot/server/scraper"
)

type serverInterface struct {
	datastore dstore.DStore
}

func New(ds dstore.DStore) api.ServerInterface {
	return serverInterface{datastore: ds}
}

// Service API
func (s serverInterface) CheckTrackedSections(ctx echo.Context) error {
	if err := mauth.VerifyAppengineCron(ctx); err != nil {
		return err
	}

	return checksections.HandleRequest(ctx, s.datastore)
}

func (s serverInterface) PruneTrackedSections(ctx echo.Context) error {
	if err := mauth.VerifyAppengineCron(ctx); err != nil {
		return err
	}

	return pruner.HandleRequest(ctx, s.datastore)
}

func (s serverInterface) ReceiveSMS(ctx echo.Context) error {
	messenger.RecieveMessage(ctx)
	return nil
}

// Public API
func (s serverInterface) GetCoursesForDepartment(ctx echo.Context, params api.GetCoursesForDepartmentParams) error {
	courses, err := s.datastore.GetCoursesForDepartment(ctx.Request().Context(), params.Term, params.DeptAbbr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
	}

	return ctx.JSON(http.StatusOK, courses)
}

func (s serverInterface) GetDepartments(ctx echo.Context) error {
	departments, err := s.datastore.GetDepartments(ctx.Request().Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
	}

	return ctx.JSON(http.StatusOK, departments)
}

func (s serverInterface) GetMeta(ctx echo.Context) error {
	meta, err := s.datastore.GetMeta(ctx.Request().Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
	}

	return ctx.JSON(http.StatusOK, *meta)
}

func (s serverInterface) GetSections(ctx echo.Context, params api.GetSectionsParams) error {
	return scraper.HandleRequest(ctx, params.Term, params.DeptAbbr, params.Course)
}

// Authenticated API
func (s serverInterface) GetUserData(ctx echo.Context, userID string) error {
	// 	_, err := mauth.VerifyToken(ctx)

	panic("implement me")
}

func (s serverInterface) UpdateUserData(ctx echo.Context, userID string) error {
	// 	_, err := mauth.VerifyToken(ctx)

	panic("implement me")
}

func (s serverInterface) RemoveTrackedSectionForUser(ctx echo.Context, userID string, sectionID string) error {
	// 	_, err := mauth.VerifyToken(ctx)
	panic("implement me")
}

func (s serverInterface) GetTrackedSectionsForUser(ctx echo.Context, userID string) error {
	// 	_, err := mauth.VerifyToken(ctx)
	panic("implement me")
}

func (s serverInterface) AddTrackedSectionsForUser(ctx echo.Context, userID string) error {
	// 	_, err := mauth.VerifyToken(ctx)
	panic("implement me")
}
