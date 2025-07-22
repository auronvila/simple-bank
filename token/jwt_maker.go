package token

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

const minSecretKeySize = 32

type JWTMaker struct {
	secretKey string
	role      string
}

func NewJwtMaker(secret string, role string) (Maker, error) {
	if len(secret) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size, must be at least %v characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey: secret, role: role}, nil
}

func (maker *JWTMaker) GenerateToken(username string, role string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, role, duration)
	if err != nil {
		return "", payload, err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	return token, payload, err
}
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, InvalidTokenError
		}
		return []byte(maker.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrorExpiredToken) {
			return nil, ErrorExpiredToken
		}
		return nil, InvalidTokenError
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, InvalidTokenError
	}
	return payload, err
}
