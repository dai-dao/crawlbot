package crawler

import (
	"net/http"
	"net/url"
)

// Request to define request structure
type Request struct {
	URL    *url.URL
	method string
}

// IResult interface
// makes it possible to compare with nil
type IResult interface {
	Request() Request
	Response() *http.Response
	Error() error
}

// Result to define network response, implements the IResult interface
// need to initialize as pointer to be recognized as IResult
type Result struct {
	req Request
	res *http.Response
	err error
}

// Request method
func (r *Result) Request() Request {
	return r.req
}

// Response method
func (r *Result) Response() *http.Response {
	return r.res
}

// Error method
func (r *Result) Error() error {
	return r.err
}
