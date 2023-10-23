package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Payload contains the payload data of the token
// Payload is our custom claims
type Payload struct {
	Number string `json:"number"`
	jwt.RegisteredClaims
}

// NewPayload creates a new token payload or claims with a specific username and duration
func NewPayload(number string, duration time.Duration) *Payload {
	payload := Payload{
		Number: number,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}

	return &payload
}
