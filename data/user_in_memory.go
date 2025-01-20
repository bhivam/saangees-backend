package data

import (
	"errors"
	"sync"
	"time"
)

type InMemoryUserStore struct {
	users      []*User
	mu         sync.Mutex
	tokenStore TokenStore
}

func NewInMemoryUserStore(tokenStore TokenStore) *InMemoryUserStore {
	return &InMemoryUserStore{
		users:      make([]*User, 0),
		tokenStore: tokenStore,
	}
}

func (userStore *InMemoryUserStore) CreateUser(user *User) (*User, error) {
	userStore.mu.Lock()
	defer userStore.mu.Unlock()

	user.ID = int64(len(userStore.users) + 1)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	userStore.users = append(userStore.users, user)
	return user, nil
}

func (userStore *InMemoryUserStore) GetUser(id int64) (*User, error) {
	userStore.mu.Lock()
	defer userStore.mu.Unlock()

	for _, user := range userStore.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (userStore *InMemoryUserStore) GetByEmail(email string) (*User, error) {
	userStore.mu.Lock()
	defer userStore.mu.Unlock()

	for _, user := range userStore.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (userStore *InMemoryUserStore) GetByToken(scope string, plaintext string) (*User, error) {
	userStore.mu.Lock()
	defer userStore.mu.Unlock()

	token, err := userStore.tokenStore.GetToken(scope, plaintext)
	if err != nil {
		return nil, err
	}

	for _, user := range userStore.users {
		if user.ID == token.UserID {
			return user, nil
		}
	}

	return nil, errors.New("user not found")
}

func (userStore *InMemoryUserStore) UpdateUser(user *User) error {
	userStore.mu.Lock()
	defer userStore.mu.Unlock()

	for i, u := range userStore.users {
		if u.ID == user.ID {
			user.UpdatedAt = time.Now()
			userStore.users[i] = user
			return nil
		}
	}
	return errors.New("user not found")
}

func (s *InMemoryUserStore) DeleteUser(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, user := range s.users {
		if user.ID == id {
			s.users = append(s.users[:i], s.users[i+1:]...)
			return nil
		}
	}
	return errors.New("user not found")
}

func (s *InMemoryUserStore) ListUsers() ([]*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return append([]*User{}, s.users...), nil
}
