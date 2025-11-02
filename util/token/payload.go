package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type Payload struct {
	Id        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expired_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	token_id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	expires := now.Add(duration)

	return &Payload{
		Id:        token_id,
		Username:  username,
		IssuedAt:  now,
		ExpiresAt: expires,
	}, nil
}

// Valid checks if token payload is valid or not
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiresAt) {
		return ErrExpiredToken
	}
	return nil
}
