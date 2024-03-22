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
		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.Error(w, "Error marshalling response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonData)
}

func (h *Handlers) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s - %s (%s)", r.Method, r.URL.Path, r.RemoteAddr)

		// Добавляем заголовки CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Link-Activation-Key")

		// Пропускаем опцион запросы
		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}
