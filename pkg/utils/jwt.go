package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTGen struct {
	SecretKey []byte
}

type jwtCustomClaims struct {
	UserId string `json:"userId"`
	jwt.StandardClaims
}

func InitJWTGen(secret string) *JWTGen {
	return &JWTGen{
		SecretKey: []byte(secret),
	}
}

func (j *JWTGen) CreateToken(userId uint) (string, error) {
	id := fmt.Sprintf("%d", userId)

	claims := &jwtCustomClaims{
		id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(10, 0, 0).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(j.SecretKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (j *JWTGen) VerifyToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return j.SecretKey, nil
	})

	if err != nil {
		return jwt.MapClaims{}, err
	}

	if !token.Valid {
		return jwt.MapClaims{}, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return jwt.MapClaims{}, fmt.Errorf("unknow probleg getting claims")
	}

	return claims, nil
}
