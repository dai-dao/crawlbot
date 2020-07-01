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

// Result to define network response
type Result struct {
	req Request
	res *http.Response
	err error
}
