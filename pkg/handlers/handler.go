package handlers

import (
	"encoding/json"
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
