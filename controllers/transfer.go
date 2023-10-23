package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/context"
	"github.com/samridhi-sahu/gobank/token"
	"github.com/samridhi-sahu/gobank/types"
	"github.com/samridhi-sahu/gobank/util"
)

type TransferRequest struct {
	FromAccount int   `json:"fromAccount"`
	ToAccount   int   `json:"toAccount"`
	Amount      int64 `json:"amount"`
}

type TransferResponse struct {
	Message string         `json:"message"`
	Account *types.Account `json:"accountDetails"`
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) {
	transferReq := TransferRequest{}
	if err := json.NewDecoder(r.Body).Decode(&transferReq); err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// taking payload from context
	authPayload := context.Get(r, authorizationPayloadKey).(*token.Payload)
	payloadNumber, err := strconv.Atoi(authPayload.Number)
	if err != nil {
		error := fmt.Errorf("while conversion in handleTransfer: %s", err.Error())
		util.WriteJSON(w, http.StatusInternalServerError, error)
		return
	}

	// check if user is transferring money form his account
	if transferReq.FromAccount != payloadNumber {
		util.WriteJSON(w, http.StatusForbidden, "you can tranfer money from your account only")
		return
	}

	// get the details of the account
	fromAccount, err := s.store.GetAccountByFilter("number", transferReq.FromAccount)
	if err != nil {
		util.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	// check balance
	if fromAccount.Balance < int64(transferReq.Amount) {
		util.WriteJSON(w, http.StatusOK, "does not have sufficient balance")
		return
	}

	// get toAccount
	toAccount, err := s.store.GetAccountByFilter("number", transferReq.ToAccount)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// debit amount from the balance of fromAccount
	prevBalance := fromAccount.Balance
	newBalance := prevBalance - transferReq.Amount
	fromAccount, err = s.store.UpdateBalance(transferReq.FromAccount, int(newBalance))
	if err != nil {
		// if error occured than abort the transaction
		s.store.UpdateBalance(transferReq.FromAccount, int(prevBalance))
		util.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	// credit amount in the balance of toAccount
	newBalance = toAccount.Balance + transferReq.Amount
	_, err = s.store.UpdateBalance(transferReq.ToAccount, int(newBalance))
	if err != nil {
		// if error occured than abort the transaction
		s.store.UpdateBalance(transferReq.FromAccount, int(prevBalance))
		util.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	transferResponse := &TransferResponse{
		Message: "Transaction Successful",
		Account: fromAccount,
	}

	err = util.WriteJSON(w, http.StatusOK, transferResponse)
	if err != nil {
		util.WriteJSON(w, http.StatusInternalServerError, err)
	}
}
