package types

import (
	"math/rand"
	"time"
)

type ApiError struct {
	Error string `json:"error"`
}

type Account struct {
	ID             int       `json:"id"`
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	HashedPassword string    `json:"-"`
	Number         int64     `json:"number"`
	Balance        int64     `json:"balance"`
	CreatedAt      time.Time `json:"createdAt"`
}

func NewAccount(firstName, lastName, hashedPassword string) *Account {
	return &Account{
		FirstName:      firstName,
		LastName:       lastName,
		HashedPassword: hashedPassword,
		Number:         int64(rand.Intn(1000000)),
		CreatedAt:      time.Now().UTC(),
	}
}
