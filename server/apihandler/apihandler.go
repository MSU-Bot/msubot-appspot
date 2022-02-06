package apihandler

import (
	"github.com/SpencerCornish/msubot-appspot/server/checksections"
	"github.com/SpencerCornish/msubot-appspot/server/dstore"
	"github.com/SpencerCornish/msubot-appspot/server/gen/api"
	"github.com/labstack/echo/v4"
)

type serverInterface struct {
	datastore dstore.DStore
}

func New(ds dstore.DStore) api.ServerInterface {
	return serverInterface{datastore: ds}
}

func (s serverInterface) CheckTrackedSections(ctx echo.Context) error {
	err := checksections.HandleRequest(ctx, s.datastore)
	return err
}

func (s serverInterface) PruneTrackedSections(ctx echo.Context) error {
	panic("implement me")
}

func (s serverInterface) GetCoursesForDepartment(ctx echo.Context, params api.GetCoursesForDepartmentParams) error {
	panic("implement me")
}

func (s serverInterface) GetDepartments(ctx echo.Context) error {
	panic("implement me")
}

func (s serverInterface) GetMeta(ctx echo.Context) error {
	panic("implement me")
}

func (s serverInterface) GetSections(ctx echo.Context, params api.GetSectionsParams) error {
	panic("implement me")
}

func (s serverInterface) ReceiveSMS(ctx echo.Context) error {
	panic("implement me")
}

func (s serverInterface) GetUserData(ctx echo.Context, userID string) error {
	panic("implement me")
}

func (s serverInterface) UpdateUserData(ctx echo.Context, userID string) error {
	panic("implement me")
}

func (s serverInterface) RemoveTrackedSectionForUser(ctx echo.Context, userID string, sectionID string) error {
	panic("implement me")
}

func (s serverInterface) GetTrackedSectionsForUser(ctx echo.Context, userID string) error {
	panic("implement me")
}

func (s serverInterface) AddTrackedSectionsForUser(ctx echo.Context, userID string) error {
	panic("implement me")
}
