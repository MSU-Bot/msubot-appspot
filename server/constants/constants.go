package constants

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

// PlivoRequest is the type sent to Plivo for texts
type PlivoRequest struct {
	Src  string `json:"src"`
	Dst  string `json:"dst"`
	Text string `json:"text"`
}

// AtlasSectionURL is the URL we hit to get section data
const AtlasSectionURL = "https://prodmyinfo.montana.edu/pls/bzagent/bzskcrse.PW_ListSchClassSimple"

// AtlasPostFormatString is the body of the request we make to Atlas for course sections. Takes three string parameters
const AtlasPostFormatString = "sel_subj=dummy&bl_online=FALSE&sel_day=dummy&term=%s&sel_subj=%s&sel_inst=ANY&sel_online=&sel_crse=%s&begin_hh=0&begin_mi=0&end_hh=0&end_mi=0"

// PlivoAPIEndpoint is the URL we hit to send SMS messages. Takes one string argument, the AuthID
const PlivoAPIEndpoint = "https://api.plivo.com/v1/Account/%s/Message/"

// PlivoSrcNum is the number MSUBot uses to send tests to users
const PlivoSrcNum = "14068000110"
