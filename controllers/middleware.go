package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/context"
	"github.com/samridhi-sahu/gobank/token"
	"github.com/samridhi-sahu/gobank/util"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// checking if token is passed in the header or not
		authorizationHeader := r.Header.Get(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			util.WriteJSON(w, http.StatusUnauthorized, err)
			return
		}

		// this will split our authorizationHeader around spce, so we should get 2 parts
		// one is Bearer and second is token
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			util.WriteJSON(w, http.StatusUnauthorized, err)
			return
		}

		// checking authorization type is bearer or not
		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			util.WriteJSON(w, http.StatusUnauthorized, err)
			return
		}

		// getting claims/payload from token
		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			util.WriteJSON(w, http.StatusUnauthorized, err)
			return
		}

		// saving payload in context, so that we can access it for authorization purpose using authorizationPayloadKey
		context.Set(r, authorizationPayloadKey, payload)

		// after authentication callind the function
		next.ServeHTTP(w, r)
	}
}
