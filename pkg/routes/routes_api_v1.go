package routes

import (
	"github.com/gorilla/mux"
)

func (r *router) initV1Routes() *mux.Router {
	r.Router.Use(r.Middleware.SetHeaders)

	r.initLinksRoutes()
	r.initUsersRoutes()
	r.initTonProofRoutes()

	return r.Router
}

func (r *router) initLinksRoutes() *mux.Router {
	linksRoutes := r.Router.PathPrefix("/links").Subrouter()
	{
		linksRoutes.HandleFunc("/create", r.Middleware.Auth(r.Handlers.CreateLink)).Methods("POST", "OPTIONS")
		linksRoutes.HandleFunc("/{slug}", r.Handlers.GetLinkBySlug).Methods("GET", "OPTIONS")
		linksRoutes.HandleFunc("/{slug}/status", r.Handlers.GetLinkStatus).Methods("GET", "OPTIONS")
		linksRoutes.HandleFunc("/{slug}/key", r.Middleware.Auth(r.Handlers.GetKeyActivation)).Methods("GET", "OPTIONS")
		linksRoutes.HandleFunc("/{slug}/activate", r.Middleware.Auth(r.Handlers.ActivateLink)).Methods("PATCH", "OPTIONS")
	}

	return r.Router
}

func (r *router) initUsersRoutes() *mux.Router {
	usersRoutes := r.Router.PathPrefix("/users").Subrouter()
	{
		usersRoutes.HandleFunc("/{address}", r.Middleware.Auth(r.Handlers.GetUserByAddress)).Methods("GET", "OPTIONS")
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

// func (r *router) initProfileRoutes() *mux.Router {
// 	profileRoutes := r.Router.PathPrefix("/profile").Subrouter()
// 	{
// 		profileRoutes.HandleFunc("/create", r.Handlers.CreateProfile).Methods("POST", "OPTIONS")
// 		profileRoutes.HandleFunc("/{id}", r.Handlers.GetProfile).Methods("GET", "OPTIONS")
// 	}
// 	return r.Router
// }

// func (r *router) initPluginRoutes() *mux.Router {
// 	pluginRoutes := r.Router.PathPrefix("/plugin").Subrouter()
// 	{
// 		pluginRoutes.HandleFunc("/{address}/generate", r.Handlers.GeneratePlugin).Methods("GET")
// 	}
// 	return r.Router
// }
