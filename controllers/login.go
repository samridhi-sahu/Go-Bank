package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/samridhi-sahu/gobank/util"
)

type loginRequest struct {
	Number   string `json:"number"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	loginRequest := loginRequest{}

	// Get the id and password from request body
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, newError("invalid request body: ", err))
		return
	}

	// Look up requested user in database
	number, err := util.ConvertIntoInt(loginRequest.Number)
	if err != nil {
		error := newError("unable to convert: ", err)
		util.WriteJSON(w, http.StatusInternalServerError, error)
		return
	}

	account, err := s.store.GetAccountByFilter("number", number)
	if err != nil {
		util.WriteJSON(w, http.StatusNotFound, err.Error())
		return
	}

	// Compare sent in password with saved user hashed password
	err = util.CheckPassword(loginRequest.Password, account.HashedPassword)
	if err != nil {
		util.WriteJSON(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Generate a JWT Token
	tokenString, err := s.tokenMaker.CreateToken(strconv.FormatInt(account.Number, 10), time.Minute*15)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	loginResponse := loginResponse{
		Token: tokenString,
	}

	// send token back in header
	err = util.WriteJSON(w, http.StatusOK, loginResponse)
	if err != nil {
		util.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func newError(message string, err error) error {
	return fmt.Errorf("%s %w", message, err)
}
