package routes

import (
	"fonates.backend/pkg/handlers"
	"fonates.backend/pkg/middlewares"
	"github.com/gorilla/mux"
)

type Router interface {
	InitRoutes(handlers handlers.Handlers) *mux.Router
}

type router struct {
	Prefix     string
	Router     *mux.Router
	Handlers   *handlers.Handlers
	Middleware *middlewares.Middleware
}

func NewRouter(prefix string) Router {
	return &router{
		Prefix:   prefix,
		Router:   mux.NewRouter().PathPrefix(prefix).Subrouter(),
		Handlers: nil,
		Middleware: nil,
	}
}

func (r *router) InitRoutes(handlers handlers.Handlers) *mux.Router {
	r.Handlers = &handlers
	r.Middleware = middlewares.NewMiddleware()

	switch r.Prefix {
	case "/api/v1":
		r.initV1Routes()
	}

	return r.Router
}
