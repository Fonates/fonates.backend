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

	var bodyLinkData models.DonationLink
	err = json.Unmarshal(body, &bodyLinkData)
	if err != nil || bodyLinkData.Name == "" {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	link := models.InitDonationLink()
	link.UserID = h.getUserId(r)
	link.Name = bodyLinkData.Name

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
		"key":    crearedLink.KeyName,
	})
}

func (h *Handlers) GetLinkBySlug(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	linkSlug := vars["slug"]

	link, err := models.InitDonationLink().GetByKey(h.Store, linkSlug)
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
	linkSlug := vars["slug"]

	if linkSlug == "" {
		h.response(w, http.StatusBadRequest, map[string]string{
			"error": "Link slug not provided",
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

	link, err := models.InitDonationLink().GetByKey(h.Store, linkSlug)
	if err != nil || link == nil {
		h.response(w, http.StatusNotFound, map[string]string{
			"error": "Link not found",
		})
		return
	}

	userId := h.getUserId(r)
	if link.UserID != userId {
		h.response(w, http.StatusForbidden, map[string]string{
			"error": "Forbidden",
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

func (h *Handlers) GetKeyActivation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	linkSlug := vars["slug"]

	link, err := models.InitDonationLink().GetByKey(h.Store, linkSlug)
	if err != nil || link == nil {
		h.response(w, http.StatusNotFound, map[string]string{
			"error": "Link not found",
		})
		return
	}

	userId := h.getUserId(r)
	if link.UserID != userId {
		h.response(w, http.StatusForbidden, map[string]string{
			"error": "Forbidden",
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

	h.response(w, http.StatusOK, key)
}

func (h *Handlers) GetLinkStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	linkSlug := vars["slug"]

	link, err := models.InitDonationLink().GetByKey(h.Store, linkSlug)
	if err != nil || link == nil {
		h.response(w, http.StatusNotFound, map[string]string{
			"error": "Link not found",
		})
		return
	}

	h.response(w, http.StatusOK, map[string]string{
		"status": link.Status,
	})
}