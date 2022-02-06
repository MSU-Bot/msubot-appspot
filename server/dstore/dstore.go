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
	GetTrackedSection(ctx context.Context, term, departmentAbbr, courseNumber string) (*models.TrackedSectionRecord, error)
	GetUser(ctx context.Context, userID string) (*models.UserRecord, error)
	GetSectionsForUser(ctx context.Context, uid string) ([]models.TrackedSectionRecord, error)
	GetAllTrackedSections(ctx context.Context) ([]models.TrackedSectionRecord, error)

	UpdateSection(ctx context.Context, sectionID string, atlasSection models.Section) error

	MoveTrackedSectionsToArchive(ctx context.Context, uids []string) error
	AddUserToExistingTrackedSection(ctx context.Context, userUID, sectionID string) error
	AddNewTrackedSection(ctx context.Context, sectionRecord models.TrackedSectionRecord) (*models.TrackedSectionRecord, error)
}
