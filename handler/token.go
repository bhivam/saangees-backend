package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/bhivam/saangees-backend/data"
	"github.com/bhivam/saangees-backend/util"
)

type TokenHandler struct {
	logger       *log.Logger
	userStore    data.UserStore
	sessionStore data.TokenStore
}

func NewTokenHandler(
	logger *log.Logger,
	userStore data.UserStore,
	sessionStore data.TokenStore,
) *TokenHandler {
	return &TokenHandler{logger, userStore, sessionStore}
}

func (tokenHandler *TokenHandler) CreateToken(
	w http.ResponseWriter,
	r *http.Request,
) {
	var requestBody CreateTokenRequestBody

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		tokenHandler.logger.Println("Error decoding request body :: ", err)
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	// TODO JSON validation

	user, err := tokenHandler.userStore.GetByEmail(requestBody.Email)
	if err != nil {
		// TODO better error handling
		tokenHandler.logger.Println("Error getting user :: ", err)
		http.Error(w, "Error getting user", http.StatusNotFound)
		return
	}

	err = util.CheckPassword(requestBody.Password, user.Hash)
	if err != nil {
		tokenHandler.logger.Println("Error checking password :: ", err)
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	token, err := tokenHandler.sessionStore.CreateToken(
		user.ID,
		24*time.Hour,
		data.ScopeAuthentication,
	)
	if err != nil {
		tokenHandler.logger.Println("Error creating token :: ", err)
		http.Error(w, "Error creating token", http.StatusInternalServerError)
		return
	}

	response := CreateTokenResponseBody{*toUserRes(user), token}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
