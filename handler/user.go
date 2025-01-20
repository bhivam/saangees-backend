package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bhivam/saangees-backend/data"
	"github.com/bhivam/saangees-backend/util"
)

type UserHandler struct {
	logger    *log.Logger
	userStore data.UserStore
}

func NewUserHandler(
	logger *log.Logger,
	userStore data.UserStore,
) *UserHandler {
	return &UserHandler{logger, userStore}
}

func (userHandler *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(util.UserContextKey{}).(*data.User)
	if !ok {
		userHandler.logger.Println("Error getting user from context")
		http.Error(w, "Error getting user from context", http.StatusInternalServerError)
		return
	}

	if !user.IsAdmin {
		userHandler.logger.Println("User is not admin")
		http.Error(w, "User is not admin", http.StatusForbidden)
		return
	}

	users, err := userHandler.userStore.ListUsers()
	if err != nil {
		userHandler.logger.Println("Error listing users :: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var res ListUsersResponseBody
	for _, user := range users {
		res = append(res, *toUserRes(user))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (userHandler *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(util.UserContextKey{}).(*data.User)
	if !ok {
		userHandler.logger.Println("Error getting user from context")
		http.Error(w, "Error getting user from context", http.StatusInternalServerError)
		return
	}

	res := *toUserRes(user)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (userHandler *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userRequestBody UserRequestBody

	err := json.NewDecoder(r.Body).Decode(&userRequestBody)
	if err != nil {
		userHandler.logger.Println("Error decoding user request body :: ", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	hashedPassword, err := util.HashPassword(userRequestBody.Password)
	if err != nil {
		userHandler.logger.Println("Error hashing password :: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	userRequestBody.Password = hashedPassword
	user := toDataUser(userRequestBody)

	created, err := userHandler.userStore.CreateUser(user)
	if err != nil {
		userHandler.logger.Println("Error creating user :: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	users, err := userHandler.userStore.ListUsers()
	if err == nil {
		json.NewEncoder(userHandler.logger.Writer()).Encode(users)
	}

	res := toUserRes(created)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
