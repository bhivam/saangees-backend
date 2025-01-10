package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/bhivam/saangees-backend/data"
	"github.com/bhivam/saangees-backend/util"
)

type SessionHandler struct {
	logger       *log.Logger
	userStore    data.UserStore
	sessionStore data.SessionStore
	tokenMaker   *util.JWTMaker
}

func NewSessionHandler(
	logger *log.Logger,
	userStore data.UserStore,
	sessionStore data.SessionStore,
	tokenMaker *util.JWTMaker,
) *SessionHandler {
	return &SessionHandler{logger, userStore, sessionStore, tokenMaker}
}

func (sessionHandler *SessionHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var loginRequestBody LoginUserRequestBody

	err := json.NewDecoder(r.Body).Decode(&loginRequestBody)
	if err != nil {
		sessionHandler.logger.Println("Error decoding user request body :: ", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	user, err := sessionHandler.userStore.GetUser(loginRequestBody.Email)
	// TODO differentiate between bad request (user not found) and server error
	if err != nil {
		sessionHandler.logger.Printf("Error getting user %v :: %v\n", loginRequestBody, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = util.CheckPassword(loginRequestBody.Password, user.Password)
	if err != nil {
		sessionHandler.logger.Println("Error checking password :: ", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	accessToken, accessClaims, err := sessionHandler.tokenMaker.GenerateToken(
		user.ID,
		user.Email,
		user.IsAdmin,
		15*time.Minute,
	)
	if err != nil {
		sessionHandler.logger.Println("Error generating access token :: ", err)
		http.Error(w, "Error creating the token", http.StatusInternalServerError)
		return
	}

	refreshToken, refreshClaims, err := sessionHandler.tokenMaker.GenerateToken(
		user.ID,
		user.Email,
		user.IsAdmin,
		30*24*time.Hour,
	)
	if err != nil {
		sessionHandler.logger.Println("Error generating refresh token :: ", err)
		http.Error(w, "Error creating the token", http.StatusInternalServerError)
		return
	}

	session, err := sessionHandler.sessionStore.CreateSession(&data.Session{
		ID:           refreshClaims.RegisteredClaims.ID,
		UserEmail:    user.Email,
		RefreshToken: refreshToken,
		IsRevoked:    false,
		ExpiresAt:    refreshClaims.RegisteredClaims.ExpiresAt.Time,
	})
	if err != nil {
		sessionHandler.logger.Println("Error creating session :: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	refreshTokenCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  refreshClaims.RegisteredClaims.ExpiresAt.Time,
		HttpOnly: true,
		Secure:   false,
		// SameSite: http.SameSiteNoneMode,
		Domain: "localhost:5173",
	}

	http.SetCookie(w, &refreshTokenCookie)

	res := LoginUserResponseBody{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessClaims.RegisteredClaims.ExpiresAt.Time,
		RefreshTokenExpiresAt: refreshClaims.RegisteredClaims.ExpiresAt.Time,
		User:                  *toUserRes(user),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (sessionHandler *SessionHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	refreshTokenCookie, err := r.Cookie("refresh_token")
	if err != nil {
		sessionHandler.logger.Println("No refresh token cookie found :: ", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	refreshClaims, err := sessionHandler.tokenMaker.VerifyToken(refreshTokenCookie.Value)
	if err != nil {
		sessionHandler.logger.Println("Error verifying refresh token :: ", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sessionID := refreshClaims.RegisteredClaims.ID

	err = sessionHandler.sessionStore.RevokedSession(sessionID)
	if err != nil {
		sessionHandler.logger.Println("Error revoking session :: ", err)
		http.Error(w, "Unable to revoke session", http.StatusInternalServerError)
		return
	}

	expireCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().AddDate(0, 0, -1),
		HttpOnly: true,
		Secure:   false,
		// SameSite: http.SameSiteNoneMode,
		Domain:   "localhost:5173",
	}

	http.SetCookie(w, &expireCookie)

	w.WriteHeader(http.StatusNoContent)
}

func (sessionHandler *SessionHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshTokenCookie, err := r.Cookie("refresh_token")
	if err != nil {
		sessionHandler.logger.Println("No refresh token cookie found :: ", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	refreshClaims, err := sessionHandler.tokenMaker.VerifyToken(refreshTokenCookie.Value)
	if err != nil {
		sessionHandler.logger.Println("Error verifying refresh token :: ", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	session, err := sessionHandler.sessionStore.GetSession(refreshClaims.RegisteredClaims.ID)
	if err != nil {
		sessionHandler.logger.Println("Error getting session :: ", err)
		http.Error(w, "Failed to get SessionID", http.StatusInternalServerError)
		return
	}

	if session.IsRevoked {
		sessionHandler.logger.Println("Session revoked")
		http.Error(w, "Session revoked", http.StatusUnauthorized)
		return
	}

	if session.UserEmail != refreshClaims.Email {
		sessionHandler.logger.Println("Session email does not match token email")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	accessToken, accessClaims, err := sessionHandler.tokenMaker.GenerateToken(
		refreshClaims.ID,
		refreshClaims.Email,
		refreshClaims.IsAdmin,
		15*time.Minute,
	)
	if err != nil {
		sessionHandler.logger.Println("Error generating access token :: ", err)
		http.Error(w, "Error creating the token", http.StatusInternalServerError)
		return
	}

	res := RefreshTokenResponseBody{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessClaims.RegisteredClaims.ExpiresAt.Time,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (sessionHandler *SessionHandler) RevokeSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["id"]
	if sessionID == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err := sessionHandler.sessionStore.RevokedSession(sessionID)
	if err != nil {
		sessionHandler.logger.Println("Error revoking session :: ", err)
		http.Error(w, "Unable to revoke session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
