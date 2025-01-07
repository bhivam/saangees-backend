package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bhivam/saangees-backend/data"
	"github.com/bhivam/saangees-backend/middleware"
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
	// TODO: Grab emails from claims

	claims := r.Context().Value(middleware.AuthKey{})
	if claims == nil {
		http.Error(w, "Error while getting data", http.StatusInternalServerError)
		return
	}

	userClaims, ok := claims.(*util.UserClaims)
	if !ok {
		http.Error(w, "Error while processing data", http.StatusInternalServerError)
		return
	}

	user, err := userHandler.userStore.GetUser(userClaims.Email)
	if err != nil {
		userHandler.logger.Println("Error retreiving user :: ", err)
		http.Error(w, "Error while getting data", http.StatusInternalServerError)
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
