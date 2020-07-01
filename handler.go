package crawler

import "net/http"

// Handler is a function used to pass into the crawler to handle the response
type Handler interface {
	Handle(Request, *http.Response, error)
}

// HandlerFunc is a function signature that implements the Handler interface
type HandlerFunc func(Request, *http.Response, error)

// Handle is the Handler interface implementation
func (h HandlerFunc) Handle(r Request, res *http.Response, err error) {
	h(r, res, err)
}
