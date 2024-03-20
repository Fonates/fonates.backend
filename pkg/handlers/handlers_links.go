package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"fonates.backend/pkg/models"
)

func (h *Handlers) CreateLinkHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var postBody models.DonationLink
	err = json.Unmarshal(body, &postBody)
	if err != nil {
		http.Error(w, "Error unmarshalling request body", http.StatusBadRequest)
		return
	}

	isValid := postBody.Validate()

	if !isValid {
		h.response(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	postBody.Status = "INACTIVE"
	crearedLink, err := postBody.Create(h.Store)
	if err != nil {
		h.response(w, http.StatusInternalServerError, map[string]string{
			"error": "Error creating link",
		})
		return
	}

	keyActivation := models.InitKeysActivation(crearedLink.ID)
	if err := keyActivation.Create(h.Store); err != nil {
		h.response(w, http.StatusInternalServerError, map[string]string{
			"error": "Error creating key activation",
		})
		return
	}

	h.response(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}