package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"gorm.io/gorm"
)

type Handlers struct {
	Store      *gorm.DB
	ServerMode string
}

func NewHandlers(store *gorm.DB, mode string) *Handlers {
	return &Handlers{
		Store:      store,
		ServerMode: mode,
	}
}

func (h *Handlers) response(w http.ResponseWriter, status int, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error marshalling response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	w.Write(jsonData)
}

func (h *Handlers) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s - %s (%s)", r.Method, r.URL.Path, r.RemoteAddr)

		w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Link-Activation-Key")
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
