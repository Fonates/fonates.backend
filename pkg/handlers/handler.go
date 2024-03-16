package handlers

import (
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

type Handlers struct {
	Store *gorm.DB
}

func NewHandlers(store *gorm.DB) *Handlers {
	return &Handlers{Store: store}
}

func (h *Handlers) response(w http.ResponseWriter, status int, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error marshalling response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonData)
}
