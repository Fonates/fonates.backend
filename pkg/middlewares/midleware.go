package middlewares

import "fonates.backend/pkg/configs"

type Middleware struct {
	SharedSecret string
}

func NewMiddleware() *Middleware {
	return &Middleware{
		SharedSecret: configs.Proof.PayloadSignatureKey,
	}
}
