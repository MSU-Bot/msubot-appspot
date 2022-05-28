package dstore

import (
	"context"

	"github.com/SpencerCornish/msubot-appspot/server/models"
)

// DStore is the interaction layer between Firestore and the go server
type DStore interface {
	GetMeta(ctx context.Context) (*models.Meta, error)
	GetDepartments(ctx context.Context) ([]string, error)
	GetCoursesForDepartment(ctx context.Context, term, department string) ([]models.DepartmentCourses, error)

	GetAllTrackedSections(ctx context.Context) ([]models.TrackedSectionRecord, error)
	GetTrackedSectionByID(ctx context.Context, ID string) (*models.TrackedSectionRecord, error)
	GetTrackedSection(ctx context.Context, term, departmentAbbr, courseNumber string) (*models.TrackedSectionRecord, error)
	GetTrackedSectionsBeforeTerm(ctx context.Context, termCondition string) ([]models.TrackedSectionRecord, error)
	GetTrackedSectionsForUser(ctx context.Context, uid string) ([]models.TrackedSectionRecord, error)

	GetUser(ctx context.Context, userID string) (*models.UserRecord, error)

	UpdateSection(ctx context.Context, sectionID string, section models.Section) error
	UpdateTrackedSection(ctx context.Context, trackedSection models.TrackedSectionRecord) error
	MoveTrackedSectionsToArchive(ctx context.Context, UIDs []string) error
	AddUserToExistingTrackedSection(ctx context.Context, userUID, sectionID string) error
	AddNewTrackedSection(ctx context.Context, sectionRecord models.TrackedSectionRecord) (*models.TrackedSectionRecord, error)
}
