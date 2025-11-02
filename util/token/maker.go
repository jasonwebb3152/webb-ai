package token

import "time"

type Maker interface {
	// Creates new token for specified user and duration
	CreateToken(username string, duration time.Duration) (string, *Payload, error)

	// Checks if token is valid or not.
	VerifyToken(token string) (*Payload, error)
}
