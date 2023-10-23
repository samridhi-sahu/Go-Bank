package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/context"
	"github.com/samridhi-sahu/gobank/token"
	"github.com/samridhi-sahu/gobank/util"
)

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) {
	id, err := util.GetId(r)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// fetch the account from db
	account, err := s.store.GetAccountByFilter("id", id)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// take the payload form context
	authPayload := context.Get(r, authorizationPayloadKey).(*token.Payload)
	payloadNumber, err := strconv.Atoi(authPayload.Number)
	if err != nil {
		error := fmt.Errorf("while conversion in handleDeleteAccount: %s", err.Error())
		util.WriteJSON(w, http.StatusInternalServerError, error)
		return
	}

	// check user is asking to delete his account or not
	if account.Number != int64(payloadNumber) {
		util.WriteJSON(w, http.StatusForbidden, "account doesn't belong to user")
		return
	}

	err = s.store.DeleteAccount(id)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	err = util.WriteJSON(w, http.StatusOK, map[string]int{"account deleted": id})
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}
}
