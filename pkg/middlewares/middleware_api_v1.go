package middlewares

import (
	"context"
	"log"
	"net/http"
	"strings"

	"fonates.backend/pkg/utils"
)

func (m *Middleware) SetHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s - %s (%s)", r.Method, r.URL.Path, r.RemoteAddr)

		w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Link-Activation-Key, Access-Control-Allow-Origin, Accpet, Authorization")
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		bearerToken := strings.Split(token, " ")
		if token == "" || len(bearerToken) != 2 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := utils.InitJWTGen(m.SharedSecret).VerifyToken(bearerToken[1])
		if err != nil || claims["userId"] == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userIdStr := claims["userId"].(string)
		ctx := context.WithValue(r.Context(), "userId", userIdStr)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
