package elements

import (
	"errors"
	"fmt"

	"github.com/boz/kcache/nsname"
)

var ErrNotFound = errors.New("Route not found")

type Route string

func NewRoute(path string) Route {
	return Route(path)
}

type Request interface {
	Route() Route
}

func NewRequest(route Route) Request {
	return request{route}
}

type request struct {
	route Route
}

func (r request) Route() Route {
	return r.route
}

type NSNameRequest interface {
	Request
	NSName() nsname.NSName
}

func NewNSNameRequest(route Route, id nsname.NSName) NSNameRequest {
	return nsNameRequest{request{route}, id}
}

type nsNameRequest struct {
	request
	id nsname.NSName
}

func (r nsNameRequest) NSName() nsname.NSName {
	return r.id
}

type Navigator interface {
	Open(Request) (Screen, error)
}

type Router interface {
	Navigator
	Register(Route, Handler)
}

type router struct {
	routes map[Route]Handler
	ctx    Context
}

func NewRouter(ctx Context) Router {
	return &router{make(map[Route]Handler), ctx}
}

func (r *router) Open(req Request) (Screen, error) {
	h, ok := r.routes[req.Route()]
	if !ok {
		return nil, ErrNotFound
	}
	return h.Open(r.ctx, req)
}

func (r *router) Register(route Route, h Handler) {
	r.routes[route] = h
}

type Handler interface {
	Open(Context, Request) (Screen, error)
}

func NewHandler(open handlerFn) Handler {
	return handler{open}
}

func NewNSNameHandler(open func(Context, NSNameRequest) (Screen, error)) Handler {
	return NewHandler(func(ctx Context, req Request) (Screen, error) {
		if req, ok := req.(NSNameRequest); ok {
			return open(ctx, req)
		}
		return nil, fmt.Errorf("Invalid request: %v expects NSName argument", req.Route())
	})
}

type handlerFn func(Context, Request) (Screen, error)

type handler struct {
	open handlerFn
}

func (h handler) Open(ctx Context, req Request) (Screen, error) {
	return h.open(ctx, req)
}
