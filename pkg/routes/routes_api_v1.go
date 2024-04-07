package routes

import (
	"fonates.backend/pkg/middlewares"
	"github.com/gorilla/mux"
)

func (r *router) initV1Routes() *mux.Router {
	r.Router.Use(middlewares.SetHeaders)

	r.initLinksRoutes()
	r.initPluginRoutes()
	r.initUsersRoutes()
	r.initTonProofRoutes()

	return r.Router
}

func (r *router) initLinksRoutes() *mux.Router {
	linksRoutes := r.Router.PathPrefix("/links").Subrouter()
	{
		linksRoutes.HandleFunc("/create", middlewares.Auth(r.Handlers.CreateLink)).Methods("POST", "OPTIONS")
		linksRoutes.HandleFunc("/{id}", r.Handlers.GetLinkByAddress).Methods("GET", "OPTIONS")
		linksRoutes.HandleFunc("/{id}/activate", r.Handlers.ActivateLink).Methods("GET", "OPTIONS")
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

func (r *router) initUsersRoutes() *mux.Router {
	usersRoutes := r.Router.PathPrefix("/users").Subrouter()
	{
		// usersRoutes.HandleFunc("/{address}/login", r.Handlers.Login).Methods("GET", "OPTIONS")
		usersRoutes.HandleFunc("/{address}/create", r.Handlers.CreateUser).Methods("POST", "OPTIONS")
	}

	return r.Router
}

func (r *router) initTonProofRoutes() *mux.Router {
	tonProofRoutes := r.Router.PathPrefix("/tonproof").Subrouter()
	{
		tonProofRoutes.HandleFunc("/generatePayload", r.Handlers.PayloadHandler).Methods("POST", "OPTIONS")
		tonProofRoutes.HandleFunc("/checkProof", r.Handlers.ProofHandler).Methods("POST", "OPTIONS")
	}

	return r.Router
}
