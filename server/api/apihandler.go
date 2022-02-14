package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/SpencerCornish/msubot-appspot/server/checksections"
	"github.com/SpencerCornish/msubot-appspot/server/dstore"
	"github.com/SpencerCornish/msubot-appspot/server/mauth"
	"github.com/SpencerCornish/msubot-appspot/server/messenger"
	"github.com/SpencerCornish/msubot-appspot/server/models"
	"github.com/SpencerCornish/msubot-appspot/server/pruner"
	"github.com/SpencerCornish/msubot-appspot/server/scraper"
	"github.com/SpencerCornish/msubot-appspot/server/usercrud"
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=types.cfg.yaml ../../api/MSUBot-Appengine-1.0.0.yaml
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=server.cfg.yaml ../../api/MSUBot-Appengine-1.0.0.yaml

type serverInterface struct {
	datastore dstore.DStore
}

func New(ds dstore.DStore) ServerInterface {
	return serverInterface{datastore: ds}
}

func sendError(ctx echo.Context, code int, message string) error {
	error := Error{
		Code:    int32(code),
		Message: message,
	}
	err := ctx.JSON(code, error)
	return err
}

// Service API
func (s serverInterface) CheckTrackedSections(ctx echo.Context) error {
	if err := mauth.VerifyAppengineCron(ctx); err != nil {
		return sendError(ctx, http.StatusForbidden, "Insufficient qualifications")
	}
	err := checksections.HandleRequest(ctx, s.datastore)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, "Server failed to handle request")
	}
	return ctx.NoContent(http.StatusOK)
}

func (s serverInterface) PruneTrackedSections(ctx echo.Context) error {
	if err := mauth.VerifyAppengineCron(ctx); err != nil {
		return sendError(ctx, http.StatusForbidden, "Insufficient qualifications")
	}
	err := pruner.HandleRequest(ctx, s.datastore)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, "Failed to prune sections")
	}

	return ctx.NoContent(http.StatusOK)
}

func (s serverInterface) ReceiveSMS(ctx echo.Context) error {
	if err := ctx.Request().ParseForm(); err != nil {
		return sendError(ctx, http.StatusBadRequest, "Formdata could not be parsed")
	}

	from := strings.Join(ctx.Request().PostForm["From"], "")
	text := strings.Join(ctx.Request().PostForm["Text"], "")
	text = strings.ToUpper(text)

	xml, err := messenger.ReceiveMessage(from, text)
	if err != nil {

		return sendError(ctx, http.StatusInternalServerError, "Failed to receive message")
	}
	return ctx.XML(http.StatusOK, xml)
}

// Public API
func (s serverInterface) GetCoursesForDepartment(ctx echo.Context, params GetCoursesForDepartmentParams) error {
	// TODO: Validate term and dept
	courseRecords, err := s.datastore.GetCoursesForDepartment(ctx.Request().Context(), params.Term, params.DeptAbbr)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, "Failed to get courses")
	}

	courses := make([]DepartmentCourse, len(courseRecords))
	for i, courseRecord := range courseRecords {
		courses[i] = DepartmentCourse{
			CourseID: courseRecord.CourseID,
			Title:    courseRecord.Title,
		}
	}

	return ctx.JSON(http.StatusOK, courses)
}

func (s serverInterface) GetDepartments(ctx echo.Context) error {
	departments, err := s.datastore.GetDepartments(ctx.Request().Context())
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, "Failed to get departments")
	}
	return ctx.JSON(http.StatusOK, departments)
}

func (s serverInterface) GetMeta(ctx echo.Context) error {
	meta, err := s.datastore.GetMeta(ctx.Request().Context())
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, "Failed to get metadata")
	}
	return ctx.JSON(http.StatusOK, *meta)
}

func (s serverInterface) GetSections(ctx echo.Context, params GetSectionsParams) error {
	scrapedSections, err := scraper.HandleRequest(ctx.Request().Context(), params.Term, params.DeptAbbr, params.Course)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, "Failed to get sections")
	}
	parsedSections, err := scrapedSectionsToParsedSections(scrapedSections)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, "Failed to parse sections")
	}

	return ctx.JSON(http.StatusOK, parsedSections)
}

// Authenticated API
func (s serverInterface) GetUserData(ctx echo.Context) error {
	token, err := mauth.VerifyToken(ctx)
	if err != nil {
		return sendError(ctx, http.StatusForbidden, "Invalid authentication")
	}

	user, err := s.datastore.GetUser(ctx.Request().Context(), token.UID)
	if err != nil {
		return sendError(ctx, http.StatusNotFound, "Failed to find user data")
	}
	respUser := User{
		UserID:      user.ID,
		Number:      user.PhoneNumber,
		WelcomeSent: &user.WelcomeSent,
		// TODO: Email: ,
	}

	return ctx.JSON(http.StatusOK, respUser)
}

func (s serverInterface) UpdateUserData(ctx echo.Context) error {
	_, err := mauth.VerifyToken(ctx)
	if err != nil {
		return sendError(ctx, http.StatusForbidden, "Invalid authentication")
	}

	panic("implement me")
}

func (s serverInterface) RemoveTrackedSectionForUser(ctx echo.Context, sectionID string) error {
	token, err := mauth.VerifyToken(ctx)
	if err != nil {
		return sendError(ctx, http.StatusForbidden, "Invalid authentication")
	}

	if err := usercrud.RemoveTrackedSection(ctx.Request().Context(), token.UID, sectionID, s.datastore); err != nil {
		return sendError(ctx, http.StatusInternalServerError, "Failed to remove section")
	}

	return ctx.NoContent(http.StatusOK)
}

func (s serverInterface) GetTrackedSectionsForUser(ctx echo.Context) error {
	_, err := mauth.VerifyToken(ctx)
	if err != nil {
		return sendError(ctx, http.StatusForbidden, "Invalid authentication")
	}

	panic("implement me")
}

func (s serverInterface) AddTrackedSectionsForUser(ctx echo.Context) error {
	_, err := mauth.VerifyToken(ctx)
	if err != nil {
		return sendError(ctx, http.StatusForbidden, "Invalid authentication")
	}

	//ADD TO API
	// const termRegex = `([0-9]){4}(30|50|70)`
	// isValidTerm, err := regexp.MatchString(termRegex, request.Term)

	panic("implement me")
}

func scrapedSectionsToParsedSections(scrapedSections []models.Section) ([]Section, error) {
	trackedSections := make([]Section, len(scrapedSections))
	for i, s := range scrapedSections {
		seats, err := strconv.Atoi(s.AvailableSeats)
		if err != nil {
			return nil, err
		}
		crn, err := strconv.Atoi(s.Crn)
		if err != nil {
			return nil, err
		}
		takenSeats, err := strconv.Atoi(s.TakenSeats)
		if err != nil {
			return nil, err
		}
		totalSeats, err := strconv.Atoi(s.TotalSeats)
		if err != nil {
			return nil, err
		}

		trackedSections[i] = Section{
			AvailableSeats: &seats,
			CourseName:     &s.CourseName,
			CourseNumber:   s.CourseNumber,
			CourseType:     &s.CourseType,
			Credits:        &s.Credits,
			Crn:            &crn,
			DeptAbbr:       s.DeptAbbr,
			Id:             nil,
			Instructor:     &s.Instructor,
			Location:       &s.Location,
			NumUsers:       nil,
			SectionNumber:  &s.SectionNumber,
			TakenSeats:     &takenSeats,
			Term:           s.Term,
			Time:           &s.Time,
			TotalSeats:     &totalSeats,
		}
	}
	return trackedSections, nil
}
