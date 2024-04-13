package handlers

import (
	"net/http"
	"strconv"

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
	userIdStr := r.Context().Value("userId")
	userId, errConvertType := strconv.Atoi(userIdStr.(string))

	if errConvertType != nil {
		h.response(w, http.StatusInternalServerError, map[string]string{
			"error": "Error converting user id",
		})
		return
	}

	if err != nil || user.ID != uint(userId) {
		h.response(w, http.StatusInternalServerError, map[string]string{
			"error": "Error getting user",
		})
		return
	}

	h.response(w, http.StatusOK, user)
}
