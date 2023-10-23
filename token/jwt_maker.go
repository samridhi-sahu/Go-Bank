package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTMaker will implement the Maker interface
type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) Maker {
	return &JWTMaker{
		secretKey: secretKey,
	}
}

// CreateToken creates a new token for a specific username and duration
func (maker *JWTMaker) CreateToken(number string, duration time.Duration) (string, error) {
	payload := NewPayload(number, duration)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	return token, err
}

// VerifyToken checks if the token is valid or not
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("token is invalid")
		}
		return []byte(maker.secretKey), nil
	}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, fmt.Errorf("invalid token: %s", err.Error())
		}
		return nil, fmt.Errorf("token is expired: %s", err.Error())
	}
	if !jwtToken.Valid {
		return nil, errors.New("invalid token")
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, errors.New("invalid token")
	}

	return payload, nil
}
