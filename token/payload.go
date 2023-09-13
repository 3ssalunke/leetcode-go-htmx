package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

type TokenPayload struct {
	ID        uuid.UUID
	Username  string
	Email     string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

func NewTokenPayload(username string, email string, duration time.Duration) (*TokenPayload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &TokenPayload{
		ID:        tokenId,
		Username:  username,
		Email:     email,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (payload *TokenPayload) Valid() error {
	if time.Now().After(payload.ExpiresAt) {
		return ErrExpiredToken
	}
	return nil
}
