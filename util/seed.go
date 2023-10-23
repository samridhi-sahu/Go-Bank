package util

import (
	"fmt"

	"github.com/samridhi-sahu/gobank/db"
	"github.com/samridhi-sahu/gobank/types"
)

// for seeding just after runing the program, pass the seed in console
// ./bin/gobank --seed
func SeedAccounts(s db.Storage) {
	seedAccount(s, "Samridhi", "Sahu", "Hello")
}

func seedAccount(store db.Storage, fname, lname, pw string) {
	account := types.NewAccount(fname, lname, pw)
	fmt.Println("new account number => ", account.Number)
	store.CreateAccount(account)
}
