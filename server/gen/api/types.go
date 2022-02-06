// Package api provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package api

// Meta defines model for Meta.
type Meta struct {

	// The total number of courses tracked by MSUBot users
	CoursesTracked int `json:"coursesTracked"`

	// The text to be used on the header of of the homepage
	Motd *string `json:"motd,omitempty"`

	// The total number of texts sent to MSUBot users
	TextsSent int `json:"textsSent"`

	// The total number of MSUBot users
	Users int `json:"users"`
}

// PlivoSMS defines model for PlivoSMS.
type PlivoSMS struct {

	// The user's phone number
	From *string `json:"From,omitempty"`

	// Set to optout if the incoming message matches with one of the standard Opt-Out keywords. Set to optin if the incoming message matches with one of the standard Opt-In keywords. Set to help if the incoming message matches with one of the standard Help keywords. Is left blank in all other cases.
	MessageIntent *string `json:"MessageIntent,omitempty"`

	// The unique identifier for the message
	MessageUUID *string `json:"MessageUUID,omitempty"`

	// The UUID of the Powerpack associated with the To phone number
	PowerpackUUID *string `json:"PowerpackUUID,omitempty"`

	// Content of the message
	Text *string `json:"Text,omitempty"`

	// Number on which the message was received
	To *string `json:"To,omitempty"`

	// Total charge for receiving the SMS (TotalRate * No. of Units)
	TotalAmount *string `json:"TotalAmount,omitempty"`

	// The charge applicable per incoming SMS unit
	TotalRate *string `json:"TotalRate,omitempty"`

	// Type of the message
	Type *string `json:"Type,omitempty"`

	// The number of parts in which the incoming message was received
	Units *int `json:"Units,omitempty"`
}

// Section defines model for Section.
type Section struct {
	AvailableSeats *int    `json:"availableSeats,omitempty"`
	CourseName     *string `json:"courseName,omitempty"`
	CourseNumber   string  `json:"courseNumber"`
	CourseType     *string `json:"courseType,omitempty"`
	Credits        *string `json:"credits,omitempty"`
	Crn            *int    `json:"crn,omitempty"`
	DeptAbbr       string  `json:"deptAbbr"`

	// MSUBot internal identifier for this section, if it exists in the system
	Id            *string `json:"id,omitempty"`
	Instructor    *string `json:"instructor,omitempty"`
	Location      *string `json:"location,omitempty"`
	NumUsers      *int    `json:"numUsers,omitempty"`
	SectionNumber *string `json:"sectionNumber,omitempty"`
	TakenSeats    *int    `json:"takenSeats,omitempty"`

	// Semester code in the format `{Year}{30|50|70|SNO}`
	Term       string  `json:"term"`
	Time       *string `json:"time,omitempty"`
	TotalSeats *int    `json:"totalSeats,omitempty"`
}

// User defines model for User.
type User struct {

	// The user's phone number
	Number string `json:"number"`

	// The user's unique identifier
	UserID string `json:"userID"`

	// whether or not the user has recieved a welcome email/text
	WelcomeSent *bool `json:"welcomeSent,omitempty"`
}

// GetCoursesForDepartmentParams defines parameters for GetCoursesForDepartment.
type GetCoursesForDepartmentParams struct {

	// Semester code in the format `{Year}{30|50|70}`
	Term string `json:"term"`

	// Short name for department
	DeptAbbr string `json:"deptAbbr"`
}

// GetSectionsParams defines parameters for GetSections.
type GetSectionsParams struct {

	// Semester code in the format `{Year}{30|50|70}`
	Term string `json:"term"`

	// Short name for department
	DeptAbbr string `json:"deptAbbr"`

	// Course "Number" (Also includes Core Tags, like RN)
	Course string `json:"course"`
}

// ReceiveSMSJSONBody defines parameters for ReceiveSMS.
type ReceiveSMSJSONBody map[string]interface{}

// UpdateUserDataJSONBody defines parameters for UpdateUserData.
type UpdateUserDataJSONBody []User

// AddTrackedSectionsForUserJSONBody defines parameters for AddTrackedSectionsForUser.
type AddTrackedSectionsForUserJSONBody []Section

// ReceiveSMSRequestBody defines body for ReceiveSMS for application/json ContentType.
type ReceiveSMSJSONRequestBody ReceiveSMSJSONBody

// UpdateUserDataRequestBody defines body for UpdateUserData for application/json ContentType.
type UpdateUserDataJSONRequestBody UpdateUserDataJSONBody

// AddTrackedSectionsForUserRequestBody defines body for AddTrackedSectionsForUser for application/json ContentType.
type AddTrackedSectionsForUserJSONRequestBody AddTrackedSectionsForUserJSONBody