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
/// [CheckTrackedSections] is run by the Appengine Cron and checks tracked CRNs
/// for open seats
func (s serverInterface) CheckTrackedSections(ctx echo.Context) error {
	if err := mauth.VerifyAppengineCron(ctx); err != nil {
		return sendError(ctx, http.StatusForbidden, "Unauthorized")
	}
	err := checksections.HandleRequest(ctx, s.datastore)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, "Server failed to check tracked sections")
	}
	return ctx.NoContent(http.StatusOK)
}

/// [PruneTrackedSections] is run by the Appengine Cron infrequently to
/// clean up sections that are irrelevant due to time
func (s serverInterface) PruneTrackedSections(ctx echo.Context) error {
	if err := mauth.VerifyAppengineCron(ctx); err != nil {
		return sendError(ctx, http.StatusForbidden, "Unauthorized")
	}
	err := pruner.HandleRequest(ctx, s.datastore)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, "Failed to prune sections")
	}

	return ctx.NoContent(http.StatusOK)
}

/// [ReceiveSMS] handles incoming SMS from Plivo, and responds with appropriate data
func (s serverInterface) ReceiveSMS(ctx echo.Context) error {
	// TODO: Validate this is coming from plivo
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
/// [GetCoursesForDepartment] Gets simple course metadata for a specified department
func (s serverInterface) GetCoursesForDepartment(ctx echo.Context, departmentID string, params GetCoursesForDepartmentParams) error {
	// TODO: Validate term and dept
	courseRecords, err := s.datastore.GetCoursesForDepartment(ctx.Request().Context(), params.Term, departmentID)
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

/// [GetDepartments] Gets the list of departments currently provided by ATLAS
/// There is no guarantee that these have classes associated with them
func (s serverInterface) GetDepartments(ctx echo.Context) error {
	departments, err := s.datastore.GetDepartments(ctx.Request().Context())
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, "Failed to get departments")
	}
	return ctx.JSON(http.StatusOK, departments)
}

/// [GetMeta] Returns any relevant global info for the client
func (s serverInterface) GetMeta(ctx echo.Context) error {
	meta, err := s.datastore.GetMeta(ctx.Request().Context())
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, "Failed to get metadata")
	}
	return ctx.JSON(http.StatusOK, *meta)
}

/// [GetSections] Scrapes ATLAS for up-to-date CRN metadata
func (s serverInterface) GetSections(ctx echo.Context, departmentID, courseID string, params GetSectionsParams) error {
	scrapedSections, err := scraper.HandleRequest(ctx.Request().Context(), params.Term, departmentID, courseID)
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
/// [GetUserData] Returns userdata for the currently auth'ed user
func (s serverInterface) GetUserData(ctx echo.Context, userID string) error {
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

/// [UpdateUserData] takes an array of userdatas and updates the currently auth'ed user
/// TODO: This takes multiple users for some reason, we probably just need to take one?
func (s serverInterface) UpdateUserData(ctx echo.Context, userID string) error {
	token, err := mauth.VerifyToken(ctx)
	if err != nil {
		return sendError(ctx, http.StatusForbidden, "Invalid authentication")
	}

	_, err = s.datastore.GetUser(ctx.Request().Context(), token.UID)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, "Failed to find user data")
	}

	// ctx.Request().Body

	panic("implement me")
}

/// [RemoveTrackedSectionForUser] takes a tracked `sectionID` and untracks it for the current user
func (s serverInterface) RemoveTrackedSectionForUser(ctx echo.Context, userID, sectionID string) error {
	token, err := mauth.VerifyToken(ctx)
	if err != nil {
		return sendError(ctx, http.StatusForbidden, "Invalid authentication")
	}

	if err := usercrud.RemoveTrackedSection(ctx.Request().Context(), token.UID, sectionID, s.datastore); err != nil {
		return sendError(ctx, http.StatusInternalServerError, "Failed to remove section")
	}

	return ctx.NoContent(http.StatusOK)
}

/// [GetTrackedSectionsForUser] Fetches the auth'ed user's currently tracked sections
func (s serverInterface) GetTrackedSectionsForUser(ctx echo.Context, userID string) error {
	_, err := mauth.VerifyToken(ctx)
	if err != nil {
		return sendError(ctx, http.StatusForbidden, "Invalid authentication")
	}

	panic("implement me")
}

/// [AddTrackedSectionsForUser] Adds new tracked sections for the auth'ed user
func (s serverInterface) AddTrackedSectionsForUser(ctx echo.Context, userID string) error {
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
