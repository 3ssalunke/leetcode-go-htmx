package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const minSecretKeySize = 32

type TokenMaker struct {
	secret string
}

func NewTokenMaker(secret string) (*TokenMaker, error) {
	if len(secret) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: key must be at least %d characters", minSecretKeySize)
	}
	return &TokenMaker{secret}, nil
}

func (maker *TokenMaker) CreateToken(username string, email string, duration time.Duration) (string, error) {
	payload, err := NewTokenPayload(username, email, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(maker.secret))
}

func (maker *TokenMaker) VerifyToken(token string) (*TokenPayload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secret), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &TokenPayload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*TokenPayload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}
