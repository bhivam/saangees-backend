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

func (tokenHandler *TokenHandler) DeleteToken(
	w http.ResponseWriter,
	r *http.Request,
) {
	user, ok := r.Context().Value(util.UserContextKey{}).(*data.User)
	if !ok {
		tokenHandler.logger.Println("Error getting user from context")
		http.Error(w, "Error getting user from token", http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	})

	if user == data.AnonymousUser {
		tokenHandler.logger.Println("Error retrieving user from context")
		http.Error(w, "User not logged in", http.StatusBadRequest)
    return
	}

	err := tokenHandler.sessionStore.DeleteAllForUser(data.ScopeAuthentication, user.ID)
	if err != nil {
		tokenHandler.logger.Println("Error deleting token :: ", err)
		http.Error(w, "Error deleting token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
	user, err := tokenHandler.userStore.GetByPhoneNumber(requestBody.PhoneNumber)
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

	response := *toUserRes(user)

	// http only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token.Plaintext,
		Expires:  token.Expiry,
		HttpOnly: true,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
