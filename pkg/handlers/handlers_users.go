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

	user, err := models.InitUser().GetByAddress(h.Store, address)
	if err != nil {
		h.response(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	h.response(w, http.StatusOK, user)
}
