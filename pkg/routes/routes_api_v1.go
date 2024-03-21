package routes

import (
	"github.com/gorilla/mux"
)

func (r *router) initV1Routes() *mux.Router {
	r.initLinksRoutes()
	r.initPluginRoutes()

	return r.Router
}

func (r *router) initLinksRoutes() *mux.Router {
	linksRoutes := r.Router.PathPrefix("/links").Subrouter()
	{
		linksRoutes.HandleFunc("/create", r.Handlers.CreateLink).Methods("POST")
		linksRoutes.HandleFunc("/{address}", r.Handlers.GetLinkByAddress).Methods("GET")
		linksRoutes.HandleFunc("/{address}/activate", r.Handlers.ActivateLink).Methods("GET")
	}

	return r.Router
}

func (r *router) initPluginRoutes() *mux.Router {
	pluginRoutes := r.Router.PathPrefix("/plugin").Subrouter()
	{
		pluginRoutes.HandleFunc("/{address}/generate", r.Handlers.GeneratePlugin).Methods("GET")
	}

	return r.Router
}
