package testutils

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// RoundTripFunc is the function type to use for request verification
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip is the wrapper to use as the http function tester
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

// MakeDummyResponse returns a dummy http response
func MakeDummyResponse() *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBufferString("DUMMYRESPONSE")),
		Header:     make(http.Header),
	}
}
