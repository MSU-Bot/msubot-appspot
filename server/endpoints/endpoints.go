package endpoints

import (
	"net/http"

	"github.com/SpencerCornish/msubot-appspot/server/apihandler"
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
	receiveSMSEndpoint = "/service/receivesms"
	healthEndpoint     = "/service/health"

	// User restricted
	getUserSectionsEndpoint       = "/user/getsections"
	trackSectionsEndpoint         = "/user/tracksections"
	removeTrackedSectionsEndpoint = "/user/removetracked"
)

func DefineServiceHandlers(handler apihandler.ApiHandler) {

	// https://msu-bot.uc.r.appspot.com/

	// Request URL: https://msu-bot.appspot.com/sections?dept=CSCI&course=112&term=202130

	// GLOBAL

	// Returns info about msubot
	http.HandleFunc(getMetaEndpoint, handler.GetMeta)

	// Returns a list of departments
	http.HandleFunc(getDepartmentsEndpoint, handler.GetDepartments)

	// Returns course names and numbers for a given department and term
	http.HandleFunc(getCoursesForDeptEndpoint, handler.GetCoursesForDept)

	// The Sections endpoint returns section metadata for a certain term, department, and course
	http.HandleFunc(getSectionsEndpoint, handler.GetSections)

	// CRON

	// the checktrackedsections is run by the cron, and is what does the actual notifying
	http.HandleFunc(checkTrackedSectionsEndpoint, handler.CheckTrackedSections)

	// The prunesections endpoint is run by the cron, and cleans up any classes from previous semesters
	// that do not need to be tracked anymore
	http.HandleFunc(pruneTrackedSectionsEndpoint, handler.PruneTrackedSections)

	// SERVICE

	// The receivemessage endpoint responds to incoming text messages
	http.HandleFunc(receiveSMSEndpoint, handler.ReceiveSMS)

	// The healthcheck endpoint reports the service's current state, as well as the latency to MSU
	http.HandleFunc(healthEndpoint, handler.HealthCheck)

	// USER

	// The Sections endpoint returns section metadata for a certain term, department, and course
	http.HandleFunc(getUserSectionsEndpoint, handler.GetUserSections)

	// The tracksections endpoint signs an authenticated user up for notifications for sections
	http.HandleFunc(trackSectionsEndpoint, handler.TrackSections)

	// The removetracked endpoint is used to remove tracked sections from a user
	http.HandleFunc(removeTrackedSectionsEndpoint, handler.RemoveTrackedSections)

}
