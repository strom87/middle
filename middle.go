package middle

import "net/http"

// Request is a regular http function
type Request func(wr http.ResponseWriter, r *http.Request)
type wrappedRequest func(wr http.ResponseWriter, r *http.Request, next Request)
type middleware func(wr http.ResponseWriter, r *http.Request) bool

type wrapper struct {
	request Request
	before  []middleware
	after   []middleware
	wrapped wrappedRequest
}

// New returns a instance of the wrapper
func New() *wrapper {
	return &wrapper{wrapped: nil}
}

// Before is used when chaining middleware for calls before the main request
func (w wrapper) Before(middlewares ...middleware) wrapper {
	w.appendMiddlewares(&w.before, middlewares)
	return w
}

// After is used when chaining middleware for calls after the main request
func (w wrapper) After(middlewares ...middleware) wrapper {
	w.appendMiddlewares(&w.after, middlewares)
	return w
}

// Wrap is used when chaining middleware for calls and encapsulates the main request
func (w wrapper) Wrap(wrapped wrappedRequest) wrapper {
	w.wrapped = wrapped
	return w
}

// UseBefore handles all the middleware calls before the main request for all the routes
func (w *wrapper) UseBefore(middlewares ...middleware) {
	w.appendMiddlewares(&w.before, middlewares)
}

// UseAfter handles all the middleware calls after the main request for all the routes
func (w *wrapper) UseAfter(middlewares ...middleware) {
	w.appendMiddlewares(&w.after, middlewares)
}

// UseWrap encapsulates the main request for all the routes
func (w *wrapper) UseWrap(wrapped wrappedRequest) {
	w.wrapped = wrapped
}

// Then runs the request and executes all the middlewares and returns a http.Handler
func (w wrapper) Then(request Request) http.Handler {
	return http.HandlerFunc(w.ThenFunc(request))
}

// ThenFunc runs the request and executes all the middlewares and returns a http function
func (w wrapper) ThenFunc(request Request) Request {
	w.request = request
	return func(wr http.ResponseWriter, r *http.Request) {
		if w.wrapped != nil {
			w.wrapped(wr, r, w.makeRequest)
		} else {
			w.makeRequest(wr, r)
		}
	}
}

func (w wrapper) makeRequest(wr http.ResponseWriter, r *http.Request) {
	if ok := w.executeMiddlewares(w.before, wr, r); !ok {
		return
	}

	w.request(wr, r)

	w.executeMiddlewares(w.after, wr, r)
}

func (w wrapper) appendMiddlewares(appender *[]middleware, middlewares []middleware) {
	for _, middleware := range middlewares {
		*appender = append(*appender, middleware)
	}
}

func (w wrapper) executeMiddlewares(middlewares []middleware, wr http.ResponseWriter, r *http.Request) bool {
	for _, middleware := range middlewares {
		if ok := middleware(wr, r); !ok {
			return false
		}
	}
	return true
}
