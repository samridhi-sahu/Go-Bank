package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/samridhi-sahu/gobank/db"
	"github.com/samridhi-sahu/gobank/token"
	"github.com/samridhi-sahu/gobank/util"
)

type APIServer struct {
	listenAddress string
	store         db.Storage
	tokenMaker    token.Maker
}

func NewAPIServer(listenAddress string, store db.Storage, tokenMaker token.Maker) *APIServer {
	return &APIServer{
		listenAddress: listenAddress,
		store:         store,
		tokenMaker:    tokenMaker,
	}
}

// Router will create new router and will handle all the routes
func (s *APIServer) SetupRouter(tokenMaker token.Maker) {
	router := mux.NewRouter()

	router.HandleFunc("/signup", s.handleCreateAccount).Methods("POST")
	router.HandleFunc("/login", s.handleLogin).Methods("POST")

	router.HandleFunc("/account", s.handleGetAccount).Methods("GET")
	router.HandleFunc("/account/add", s.addBalance).Methods("POST")

	router.HandleFunc("/account/{id}", authMiddleware(tokenMaker, s.handleGetAccountByID)).Methods("GET")
	router.HandleFunc("/account/{id}", authMiddleware(tokenMaker, s.handleDeleteAccount)).Methods("DELETE")
	router.HandleFunc("/transfer", authMiddleware(tokenMaker, s.handleTransfer)).Methods("POST")

	log.Println("JSON API server running on port: ", s.listenAddress)

	http.ListenAndServe(s.listenAddress, router)
}

// route to add balance
type AddRequest struct {
	Number int `json:"number"`
	Amount int `json:"amount"`
}

func (s *APIServer) addBalance(w http.ResponseWriter, r *http.Request) {
	addReq := AddRequest{}
	if err := json.NewDecoder(r.Body).Decode(&addReq); err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// get the details of the account
	account, err := s.store.GetAccountByFilter("number", addReq.Number)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// check amount
	if addReq.Amount < 0 {
		util.WriteJSON(w, http.StatusNotFound, "negative amount does not allowed")
		return
	}

	// add balance
	newBalance := int(account.Balance) + addReq.Amount
	_, err = s.store.UpdateBalance(addReq.Number, int(newBalance))
	if err != nil {
		util.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	message := fmt.Sprintf("Updated Amount: %d", newBalance)
	err = util.WriteJSON(w, http.StatusOK, message)
	if err != nil {
		util.WriteJSON(w, http.StatusInternalServerError, err)
	}
}
