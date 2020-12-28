package models

import "time"

// Section is a section model
type Section struct {
	Term string

	DeptAbbr string
	DeptName string

	CourseName   string
	CourseNumber string
	CourseType   string
	Credits      string

	Instructor string
	Time       string
	Location   string

	SectionNumber  string
	Crn            string
	TotalSeats     string
	TakenSeats     string
	AvailableSeats string
}

// TrackedSectionRecord is a datastore record for tracked sections in firestore
type TrackedSectionRecord struct {
	ID string

	CourseName     string    `firestore:"courseName"`
	CourseNumber   string    `firestore:"courseNumber"`
	CreationTime   time.Time `firestore:"creationTime"`
	Crn            string    `firestore:"crn"`
	Department     string    `firestore:"department"`
	DepartmentAbbr string    `firestore:"departmentAbbr"`
	Instructor     string    `firestore:"instructor"`
	OpenSeats      string    `firestore:"openSeats"`
	SectionNumber  int       `firestore:"sectionNumber"`
	Term           string    `firestore:"term"`
	TotalSeats     string    `firestore:"totalSeats"`
	Users          []string  `firestore:"users"`
}

// Meta is information about the service
type Meta struct {
	CoursesTracked int    `firestore:"coursesTracked"`
	Users          int    `firestore:"users"`
	TextsSent      int    `firestore:"textsSent"`
	Motd           string `firestore:"motd"`
}

type DepartmentCourses struct {
	CourseID string
	Title    string
}
