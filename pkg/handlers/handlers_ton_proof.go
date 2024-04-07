package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"fonates.backend/pkg/ton"
	"github.com/golang-jwt/jwt"
	"github.com/tonkeeper/tongo"
)

type jwtCustomClaims struct {
	Address string `json:"address"`
	jwt.StandardClaims
}

// type handler struct {
// sharedSecret string
// payloadTtl   time.Duration
// }

// func newHandler(sharedSecret string, payloadTtl time.Duration) *handler {
// 	h := handler{
// 		sharedSecret: sharedSecret,
// 		payloadTtl:   payloadTtl,
// 	}
// 	return &h
// }

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

	fmt.Println(tp.Proof.Payload)
	fmt.Println(h.SharedSecret)
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

	claims := &jwtCustomClaims{
		tp.Address,
		jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(10, 0, 0).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(h.SharedSecret))
	if err != nil {
		h.response(w, http.StatusInternalServerError, map[string]string{
			"error": "Error signing token",
		})
	}

	h.response(w, http.StatusOK, map[string]interface{}{
		"token": t,
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
