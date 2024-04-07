package handlers

import (
	"net/http"

	"fonates.backend/pkg/models"
	"fonates.backend/pkg/utils"
	"github.com/gorilla/mux"
)

func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	addr := vars["address"]

	if addr == "" || !utils.ValidateTonAddress(addr) {
		h.response(w, http.StatusBadRequest, map[string]string{
			"error": "Address not provided",
		})
		return
	}

	user := models.InitUser()
	user.Address = addr

	// createdUser, err := user.Create(h.Store)
	// if err != nil {
	// 	h.response(w, http.StatusInternalServerError, map[string]string{
	// 		"error": "Error creating user",
	// 	})
	// 	return
	// }

	// h.response(w, http.StatusOK, map[string]interface{}{
	// 	"token": token,
	// })
}
