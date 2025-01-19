package auth

import (
	"context"
	"strconv"

	"firebase.google.com/go/v4/auth"
	"golang.org/x/crypto/bcrypt"

	"github.com/bhivam/saangees-backend/data"
)

type AuthService struct {
	userStore data.UserStore
	fireAuth  *auth.Client
}

func NewAuthService(userStore data.UserStore, fireAuth *auth.Client) *AuthService {
	return &AuthService{
		userStore: userStore,
		fireAuth:  fireAuth,
	}
}

func (s *AuthService) Login(email string, password string) (string, error) {
	user, err := s.userStore.GetUser(email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", err
	}

	token, err := s.fireAuth.CustomToken(context.Background(), strconv.FormatInt(user.ID, 10))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Register(email string, password string, isAdmin bool, name string) (string, error) {
	_, err := s.userStore.GetUser(email)

	if err != (&data.UserNotFoundError{}) {
		return "", err
	}
  
  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
  if err != nil {
    return "", err
  }

  user := &data.User{
    Name: name,
    Email: email,
    Password: string(hashedPassword),
    IsAdmin: isAdmin,
  }
  
  user, err = s.userStore.CreateUser(user)
  if err != nil {
    return "", err
  }

  token, err := s.fireAuth.CustomToken(context.Background(), strconv.FormatInt(user.ID, 10))
  if err != nil {
    return "", err
  }

  return token, nil
}
