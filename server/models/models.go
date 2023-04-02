package models

// Section is a section model
type Section struct {
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

type CourseMeta struct {
	Number string
	Title  string
}

type Department struct {
	Id   string
	Name string
}
