package datastore

import "github.com/SpencerCornish/msubot-appspot/server/models"

// DataStore is the interaction layer between Firestore and the go server
type DataStore interface {
	GetMeta() (models.Meta, error)
	GetDepartments() ([]string, error)
	GetCoursesForDepartment() ([]string, error)

	GetSectionsForUser(uid string) ([]models.TrackedSectionRecord, error)
	GetAllTrackedSections() ([]models.TrackedSectionRecord, error)

	MoveTrackedSectionsToArchive(uids []string) error
	AddUserToTrackedSection(userUID, sectionUID string) error
}
