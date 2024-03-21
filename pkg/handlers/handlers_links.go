package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"fonates.backend/pkg/models"
	"github.com/gorilla/mux"
)

func (h *Handlers) CreateLink(w http.ResponseWriter, r *http.Request) {
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

func (h *Handlers) GetLinkByAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	link, err := models.InitDonationLink().GetByAddress(h.Store, address)
	if err != nil {
		h.response(w, http.StatusNotFound, map[string]string{
			"error": "Link not found",
		})
		return
	}

	if link.Status == "INACTIVE" {
		h.response(w, http.StatusNotFound, map[string]string{
			"error": "Link not activated",
		})
		return
	}

	h.response(w, http.StatusOK, link)
}

func (h *Handlers) ActivateLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	if address == "" {
		h.response(w, http.StatusBadRequest, map[string]string{
			"error": "Address not provided",
		})
		return
	}

	keyActivation := r.Header.Get("X-Link-Activation-Key")
	if keyActivation == "" {
		h.response(w, http.StatusBadRequest, map[string]string{
			"error": "Key not provided",
		})
		return
	}

	link, err := models.InitDonationLink().GetByAddress(h.Store, address)
	if err != nil {
		h.response(w, http.StatusNotFound, map[string]string{
			"error": "Link not found",
		})
		return
	}

	key, err := models.InitKeysActivation(link.ID).GetByLinkID(h.Store, link.ID)
	if err != nil {
		h.response(w, http.StatusNotFound, map[string]string{
			"error": "Key not found",
		})
		return
	}

	if key.Status == "ACTIVE" {
		h.response(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
		return
	}

	if key.Key.String() != keyActivation {
		h.response(w, http.StatusForbidden, map[string]string{
			"error": "Invalid key",
		})
		return
	}

	if err := link.Activate(h.Store); err != nil {
		h.response(w, http.StatusInternalServerError, map[string]string{
			"error":  "Error activating link",
		})
		return
	}

	if err := key.Activate(h.Store); err != nil {
		h.response(w, http.StatusInternalServerError, map[string]string{
			"error": "Error activating key",
		})
		return
	}

	h.response(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}
