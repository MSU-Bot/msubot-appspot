package serverutils

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/MSU-Bot/msubot-appspot/server/testutils"
	"github.com/stretchr/testify/assert"
)

const (
	testTerm           = "909090"
	testDept           = "YEET"
	testCourse         = "290IN"
	testPlivoID        = "PlivoID"
	testPlivoAuthToken = "PlivoAuthToken"
	testNumber         = "+15559992222"
	testTextContent    = "Hello Test 1234"
)

func TestMakeAtlasSectionRequest_Success(t *testing.T) {

	client := testutils.NewTestClient(func(req *http.Request) *http.Response {
		assert.Equal(t, sectionRequestURL, req.URL.String())

		expectedReqBodyReader := buildAtlasRequestBody(testTerm, testDept, testCourse)

		buf := new(bytes.Buffer)
		buf.ReadFrom(expectedReqBodyReader)
		expectedReqBody := buf.String()

		buf = new(bytes.Buffer)
		buf.ReadFrom(req.Body)
		actualReqBody := buf.String()

		assert.Equal(t, expectedReqBody, actualReqBody)

		return testutils.MakeDummyResponse()
	})

	response, err := MakeAtlasSectionRequest(client, testTerm, testDept, testCourse)
	assert.Nil(t, err)
	assert.Equal(t, testutils.MakeDummyResponse(), response)
}

func TestBuildAtlasRequestBody_Success(t *testing.T) {
	expectedOutput := "sel_subj=dummy&bl_online=FALSE&sel_day=dummy&sel_online=dummy&term=909090&sel_subj=YEET&sel_inst=0&sel_online=&sel_crse=290IN&begin_hh=0&begin_mi=0&end_hh=0&end_mi=0"
	actual := buildAtlasRequestBody(testTerm, testDept, testCourse)

	buf := new(bytes.Buffer)
	buf.ReadFrom(actual)
	actualString := buf.String()
	assert.Equal(t, expectedOutput, actualString)
}

// func TestSendText_SingleSuccess(t *testing.T) {

// 	os.Setenv("PLIVO_AUTH_ID", testPlivoID)
// 	os.Setenv("PLIVO_AUTH_TOKEN", testPlivoAuthToken)

// 	client := testutils.NewTestClient(func(req *http.Request) *http.Response {
// 		assert.Equal(t, getPlivoURL(testPlivoID), req.URL.String())

// 		username, pass, ok := req.BasicAuth()
// 		assert.True(t, ok)
// 		assert.Equal(t, testPlivoID, username)
// 		assert.Equal(t, testPlivoAuthToken, pass)

// 		assert.Equal(t, "POST", req.Method)

// 		requestBuffer := new(bytes.Buffer)
// 		requestBuffer.ReadFrom(req.Body)
// 		expectedRequest := &plivoRequest{
// 			Src:  sourceNumber,
// 			Dst:  testNumber,
// 			Text: testTextContent,
// 		}
// 		var actualRequest plivoRequest
// 		json.Unmarshal(requestBuffer.Bytes(), &actualRequest)

// 		assert.Equal(t, expectedRequest, &actualRequest)

// 		return testutils.MakeDummyResponse()
// 	})

// 	response, err := SendText(client, testNumber, testTextContent)
// 	assert.Nil(t, err)
// 	assert.Equal(t, testutils.MakeDummyResponse(), response)
// }
