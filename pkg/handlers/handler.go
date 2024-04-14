package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"fonates.backend/pkg/configs"
	"gorm.io/gorm"
)

type Handlers struct {
	Store        *gorm.DB
	ServerMode   string
	SharedSecret string
	PayloadTtl   time.Duration
}

func NewHandlers(store *gorm.DB, mode string) *Handlers {
	return &Handlers{
		Store:        store,
		ServerMode:   mode,
		SharedSecret: configs.Proof.PayloadSignatureKey,
		PayloadTtl:   time.Duration(configs.Proof.ProofLifeTimeSec) * time.Second,
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

func (h *Handlers) getUserId(r *http.Request) uint {
	userIdStr := r.Context().Value("userId")
	userId, err := strconv.Atoi(userIdStr.(string))
	if err != nil {
		return 0
	}

	return uint(userId)
}
