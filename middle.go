package middle

import "net/http"

// Request is a regular http function
type Request func(wr http.ResponseWriter, r *http.Request)
type wrappedRequest func(wr http.ResponseWriter, r *http.Request, next Request)
type middleware func(wr http.ResponseWriter, r *http.Request) bool

type wrapper struct {
	before     []middleware
	after      []middleware
	wrapped    wrappedRequest
	useBefore  []middleware
	useAfter   []middleware
	useWrapped wrappedRequest
}

// New returns a instance of the wrapper
func New() *wrapper {
	return &wrapper{wrapped: nil, useWrapped: nil}
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
	w.appendMiddlewares(&w.useBefore, middlewares)
}

// UseAfter handles all the middleware calls after the main request for all the routes
func (w *wrapper) UseAfter(middlewares ...middleware) {
	w.appendMiddlewares(&w.useAfter, middlewares)
}

// UseWrap encapsulates the main request for all the routes
func (w *wrapper) UseWrap(wrapped wrappedRequest) {
	w.useWrapped = wrapped
}

// Then runs the request and executes all the middlewares
func (w wrapper) Then(req Request) Request {
	return func(wr http.ResponseWriter, r *http.Request) {
		if ok := w.executeMiddlewares(w.useBefore, wr, r); !ok {
			return
		} else if ok := w.executeMiddlewares(w.before, wr, r); !ok {
			return
		}

		if w.wrapped != nil {
			w.wrapped(wr, r, req)
		} else if w.useWrapped != nil {
			w.useWrapped(wr, r, req)
		} else {
			req(wr, r)
		}

		if ok := w.executeMiddlewares(w.useAfter, wr, r); !ok {
			return
		}
		w.executeMiddlewares(w.after, wr, r)
	}
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
