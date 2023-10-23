package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// WriteJSON will set all the details in the ResponseWriter
func WriteJSON(w http.ResponseWriter, status int, value any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	error := json.NewEncoder(w).Encode(value)
	return error
}

// getId will return the id passed in the url as integer
func GetId(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}
	return id, nil
}

func ConvertIntoInt(s string) (int, error) {
	val, err := strconv.Atoi(s)
	if err != nil {
		return val, fmt.Errorf("invalid string %s", s)
	}
	return val, nil
}
