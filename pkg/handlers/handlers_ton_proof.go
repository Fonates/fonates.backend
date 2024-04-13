package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"fonates.backend/pkg/models"
	"fonates.backend/pkg/ton"
	"fonates.backend/pkg/utils"
	"github.com/tonkeeper/tongo"
)

func (h *Handlers) ProofHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var tp ton.TonProof
	err = json.Unmarshal(body, &tp)
	if err != nil {
		h.response(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	err = ton.CheckPayload(tp.Proof.Payload, h.SharedSecret)
	if err != nil {
		h.response(w, http.StatusBadRequest, map[string]string{
			"error": "payload verification failed: " + err.Error(),
		})
		return
	}

	parsed, err := ton.ConvertTonProofMessage(&tp)
	if err != nil {
		h.response(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	net := ton.TonNetworks[tp.Network]
	if net == nil {
		h.response(w, http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("invalid network: %v", tp.Network),
		})
		return
	}

	addr, err := tongo.ParseAccountID(tp.Address)
	if err != nil {
		h.response(w, http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("invalid account: %v", tp.Address),
		})
		return
	}

	ctx := r.Context()
	check, err := ton.CheckProof(ctx, addr, net, parsed)
	if err != nil {
		h.response(w, http.StatusBadRequest, map[string]string{
			"error": "proof checking error: " + err.Error(),
		})
		return
	}

	if !check {
		h.response(w, http.StatusBadRequest, map[string]string{
			"error": "proof verification failed",
		})
		return
	}

	addrHuman, errConvert := ton.AddrFriendly(tp.Address, tp.Network)
	if errConvert != nil {
		log.Println("Error converting address: ", errConvert)
		h.response(w, http.StatusInternalServerError, map[string]string{
			"error": "Error converting address",
		})
		return
	}

	user := models.InitUser()
	foundUser, errFound := user.GetByAddress(h.Store, addrHuman)
	if errFound != nil {
		h.response(w, http.StatusInternalServerError, map[string]string{
			"error": "Error getting user",
		})
		return
	}

	if foundUser.ID == 0 {
		user.Address = addrHuman
		foundUser, err = user.Create(h.Store)
		if err != nil {
			h.response(w, http.StatusInternalServerError, map[string]string{
				"error": "Error creating user",
			})
			return
		}
	}

	jwtToken, errGenerate := utils.InitJWTGen(h.SharedSecret).CreateToken(foundUser.ID)
	if errGenerate != nil {
		h.response(w, http.StatusInternalServerError, map[string]string{
			"error": "Error signing token",
		})
	}

	h.response(w, http.StatusOK, map[string]interface{}{
		"token": jwtToken,
	})
}

func (h *Handlers) PayloadHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := ton.GeneratePayload(h.SharedSecret, h.PayloadTtl)
	if err != nil {
		h.response(w, http.StatusInternalServerError, map[string]string{
			"error": "Error generating payload",
		})
		return
	}

	h.response(w, http.StatusOK, map[string]string{
		"payload": payload,
	})
}
