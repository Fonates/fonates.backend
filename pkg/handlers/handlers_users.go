package handlers

import (
	"net/http"

	"fonates.backend/pkg/models"
	"github.com/gorilla/mux"
)

func (h *Handlers) GetUserByAddress(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	address := params["address"]

	if address == "" {
		h.response(w, http.StatusBadRequest, map[string]string{
			"error": "address is required",
		})
		return
	}

	ctxUserId := h.getUserId(r)
	user, err := models.InitUser().GetByAddress(h.Store, address)

	if err != nil || ctxUserId == 0 || user.ID == 0 {
		h.response(w, http.StatusInternalServerError, map[string]string{
			"error": "Error getting user",
		})
		return
	}

	h.response(w, http.StatusOK, user)
}
