package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/samridhi-sahu/gobank/types"
	"github.com/samridhi-sahu/gobank/util"
)

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	createAccountReq := CreateAccountRequest{}
	if err := json.NewDecoder(r.Body).Decode(&createAccountReq); err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// hash the password
	hashedPassword, err := util.HashedPassword(createAccountReq.Password)
	if err != nil {
		util.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	// create new account
	account := types.NewAccount(createAccountReq.FirstName, createAccountReq.LastName, hashedPassword)
	id, err := s.store.CreateAccount(account)
	if err != nil {
		util.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	// set the id of account
	account.ID = id

	// response
	err = util.WriteJSON(w, http.StatusOK, account)
	if err != nil {
		util.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
}
