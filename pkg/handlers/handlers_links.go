package handlers

import (
	"log"
	"net/http"

	"fonates.backend/pkg/models"
	"github.com/gorilla/mux"
)

func (h *Handlers) CreateLink(w http.ResponseWriter, r *http.Request) {
	link := models.InitDonationLink()
	link.UserID = r.Context().Value("userId").(uint)

	log.Printf("User ID: %v", link.UserID)

	crearedLink, err := link.Create(h.Store)
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
	linkId := vars["id"]

	link, err := models.InitDonationLink().GetById(h.Store, linkId)
	if err != nil || link == nil {
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
	linkId := vars["id"]

	if linkId == "" {
		h.response(w, http.StatusBadRequest, map[string]string{
			"error": "Id not provided",
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

	link, err := models.InitDonationLink().GetById(h.Store, linkId)
	if err != nil || link == nil {
		h.response(w, http.StatusNotFound, map[string]string{
			"error": "Link not found",
		})
		return
	}

	key, err := models.InitKeysActivation(link.ID).GetByLinkID(h.Store)
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
			"error": "Error activating link",
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
