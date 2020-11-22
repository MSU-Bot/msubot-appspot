package endpoints

import (
	"net/http"

	"github.com/SpencerCornish/msubot-appspot/server/checksections"
	"github.com/SpencerCornish/msubot-appspot/server/healthcheck"
	"github.com/SpencerCornish/msubot-appspot/server/messenger"
	"github.com/SpencerCornish/msubot-appspot/server/pruner"
	"github.com/SpencerCornish/msubot-appspot/server/removetracked"
	"github.com/SpencerCornish/msubot-appspot/server/scraper"
	"github.com/SpencerCornish/msubot-appspot/server/tracksections"
)

const (
	// Section interaction

	// Global
	getMetaEndpoint           = "/global/getmeta"
	getDepartmentsEndpoint    = "/global/getdepartments"
	getCoursesForDeptEndpoint = "/global/coursesfordept"
	getSectionsEndpoint       = "/global/getsections"

	// Cron restricted
	checkTrackedSectionsEndpoint = "/cron/checktrackedsections"
	pruneTrackedSectionsEndpoint = "/cron/prunetrackedsections"

	// Service communications
	recieveSMSEndpoint = "/service/recievesms"
	healthEndpoint     = "/service/health"

	// User restricted
	getUserSectionsEndpoint       = "/user/getsections"
	trackSectionsEndpoint         = "/user/tracksections"
	removeTrackedSectionsEndpoint = "/user/removetracked"
)

func DefineServiceHandlers() {

	// GLOBAL

	// Returns info about msubot
	http.HandleFunc(getMetaEndpoint, meta.GetMeta)

	// Returns a list of departments
	http.HandleFunc(getDepartmentsEndpoint, meta.GetDepartments)

	// Returns course names and numbers for a given department and term
	http.HandleFunc(getCoursesForDeptEndpoint, meta.GetCoursesForDept)

	// The Sections endpoint returns section metadata for a certain term, department, and course
	http.HandleFunc(getSectionsEndpoint, scraper.HandleRequest)

	// CRON

	// the checktrackedsections is run by the cron, and is what does the actual notifying
	http.HandleFunc(checkTrackedSectionsEndpoint, checksections.HandleRequest)

	// The prunesections endpoint is run by the cron, and cleans up any classes from previous semesters
	// that do not need to be tracked anymore
	http.HandleFunc(pruneTrackedSectionsEndpoint, pruner.HandleRequest)

	// SERVICE

	// The recievemessage endpoint responds to incoming text messages
	http.HandleFunc(recieveSMSEndpoint, messenger.RecieveMessage)

	// The healthcheck endpoint reports the service's current state, as well as the latency to MSU
	http.HandleFunc(healthEndpoint, healthcheck.CheckHealth)

	// USER

	// The Sections endpoint returns section metadata for a certain term, department, and course
	http.HandleFunc(getUserSectionsEndpoint, scraper.HandleRequest)

	// The tracksections endpoint signs an authenticated user up for notifications for sections
	http.HandleFunc(trackSectionsEndpoint, tracksections.HandleRequest)

	// The removetracked endpoint is used to remove tracked sections from a user
	http.HandleFunc(removeTrackedSectionsEndpoint, removetracked.HandleRequest)

}

func defineInternalHandlers() {

}
