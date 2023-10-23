package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/context"
	"github.com/samridhi-sahu/gobank/token"
	"github.com/samridhi-sahu/gobank/util"
)

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) {
	accounts, err := s.store.GetAllAccounts()
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := util.WriteJSON(w, http.StatusOK, accounts); err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}
}

func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) {
	id, err := util.GetId(r)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// fet the account from db
	account, err := s.store.GetAccountByFilter("id", id)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// fetch the payload from context
	authPayload := context.Get(r, authorizationPayloadKey).(*token.Payload)
	payloadNumber, err := strconv.Atoi(authPayload.Number)
	if err != nil {
		error := fmt.Errorf("while conversion in handleAccountGetByID: %s", err.Error())
		util.WriteJSON(w, http.StatusInternalServerError, error)
		return
	}

	// check if user is asking for his data or not
	if account.Number != int64(payloadNumber) {
		util.WriteJSON(w, http.StatusForbidden, "account doesn't belong to user")
		return
	}

	err = util.WriteJSON(w, http.StatusOK, account)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}
}
