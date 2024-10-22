package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"fonates.backend/pkg/models"
	"github.com/gorilla/mux"
)

const (
	ALERT_STATUS_WAIT   = "WAIT"
	ALERT_STATUS_SENDED = "SENDED"
)

type Alert struct {
	Data   interface{}
	Status string
}

var dataChan = make(map[string]chan Alert)

func (h *Handlers) CreateDonate(w http.ResponseWriter, r *http.Request) {
	body, errReadBody := io.ReadAll(r.Body)
	if errReadBody != nil {
		h.response(w, http.StatusInternalServerError, map[string]string{
			"error": "Error reading request body",
		})
		return
	}

	var donate = models.InitDonate()
	if err := json.Unmarshal(body, &donate); err != nil {
		h.response(w, http.StatusInternalServerError, map[string]string{
			"error": "Error unmarshal request body",
		})
		return
	}

	if donate.Amount == 0 || donate.Username == "" || donate.DonationLinkID == 0 || donate.Hash == "" {
		h.response(w, http.StatusBadRequest, map[string]string{
			"error": "Amount, Username, DonationLinkID, Hash are required",
		})
		return
	}

	userId := h.getUserId(r)
	if userId == 0 {
		h.response(w, http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
		return
	}

	donationLink := models.InitDonationLink()
	donationLink.ID = donate.DonationLinkID
	donate.UserID = userId

	if errFound := donationLink.GetById(h.Store); errFound != nil {
		h.response(w, http.StatusInternalServerError, map[string]string{
			"error": "Error getting donation link",
		})
		return
	}

	if donationLink.Status != models.LINK_ACTIVE || donationLink.Status == models.LINK_BLOCKED {
		h.response(w, http.StatusBadRequest, map[string]string{
			"error": "Donation link is not active",
		})
		return
	}

	if err := donate.Create(h.Store); err != nil {
		h.response(w, http.StatusInternalServerError, map[string]string{
			"error": "Error creating donate",
		})
		return
	}

	if _, ok := dataChan[donationLink.KeyName]; !ok {
		dataChan[donationLink.KeyName] = make(chan Alert)
	}

	go func(c chan Alert) {
		c <- Alert{
			Status: ALERT_STATUS_WAIT,
			Data:   donate,
		}
	}(dataChan[donationLink.KeyName])

	h.response(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func (h *Handlers) DonatesStreaming(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	link_key := vars["link_key"]

	if link_key == "" {
		h.response(w, http.StatusBadRequest, map[string]string{
			"error": "Link ID is required",
		})
		return
	}

	donationLink, err := models.InitDonationLink().GetByKey(h.Store, link_key)
	if err != nil || donationLink == nil {
		h.response(w, http.StatusInternalServerError, map[string]string{
			"error": "Error getting donation link",
		})
		return
	}

	log.Default().Print("DONATION LINK: ", donationLink)

	if donationLink.Status != models.LINK_ACTIVE || donationLink.Status == models.LINK_BLOCKED {
		h.response(w, http.StatusBadRequest, map[string]string{
			"error": "Donation link is not active",
		})
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	fmt.Fprintf(w, "data: %s\n\n", "heartbeat")

	var ticker = time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	for {
		select {
		case alert := <-dataChan[link_key]:
			log.Default().Print("GET CHANNEL: ", alert)
			if alert.Status != ALERT_STATUS_WAIT {
				continue
			}

			jsonData, err := json.Marshal(alert.Data)
			if err != nil {
				log.Println(err)
				continue
			}

			delete(dataChan, link_key)

			fmt.Fprintf(w, "data: %s\n\n", jsonData)
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}
		case <-ticker.C:
			fmt.Fprintf(w, "data: %s\n\n", "heartbeat")

			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}

			select {
			case <-r.Context().Done():
				return
			default:
			}
		case <-r.Context().Done():
			return
		}
	}
}
