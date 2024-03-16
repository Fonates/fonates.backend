package routes

import (
	"net/http"

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
		linksRoutes.HandleFunc("/create", r.Handlers.CreateLinkHandler).Methods("POST")
		linksRoutes.HandleFunc("/{address}", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")
		linksRoutes.HandleFunc("/{address}/activate", func(w http.ResponseWriter, r *http.Request) {}).Methods("UPDATE")
	}

	return r.Router
}

func (r *router) initPluginRoutes() *mux.Router {
	pluginRoutes := r.Router.PathPrefix("/plugins").Subrouter()
	{
		pluginRoutes.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {}).Methods("POST")
		pluginRoutes.HandleFunc("/{address}", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")
		pluginRoutes.HandleFunc("/{address}/activate", func(w http.ResponseWriter, r *http.Request) {}).Methods("UPDATE")
	}

	return r.Router
}
