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
	CourseName     string
	CourseNumber   string
	CreationTime   time.Time
	Crn            string
	Department     string
	DepartmentAbbr string
	Instructor     string
	OpenSeats      string
	SectionNumber  int
	Term           string
	TotalSeats     string
	Users          []string
}

// Meta is information about the service
type Meta struct {
	TotalTracked    int
	TotalUsers      int
	MessageOfTheDay string
}
