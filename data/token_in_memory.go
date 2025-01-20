package data

import (
	"bytes"
	"crypto/sha256"
	"sync"
	"time"
)

type InMemoryTokenStore struct {
	tokens []*Token
	mu     sync.Mutex
}

func NewInMemoryTokenStore() *InMemoryTokenStore {
	return &InMemoryTokenStore{
		tokens: []*Token{},
	}
}

func (store *InMemoryTokenStore) CreateToken(
	userID int64,
	ttl time.Duration,
	scope string,
) (*Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	store.tokens = append(store.tokens, token)
	return token, err
}

func (store *InMemoryTokenStore) InsertToken(token *Token) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	store.tokens = append(store.tokens, token)
	return nil
}

func (store *InMemoryTokenStore) DeleteAllForUser(scope string, userID int64) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	var remainingTokens []*Token
	for _, token := range store.tokens {
		if !(token.Scope == scope && token.UserID == userID) {
			remainingTokens = append(remainingTokens, token)
		}
	}

	store.tokens = remainingTokens
	return nil
}

func (store *InMemoryTokenStore) GetToken(scope string, plaintext string) (*Token, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	hash := sha256.Sum256([]byte(plaintext))

	for _, token := range store.tokens {
		if token.Scope == scope && bytes.Equal(token.Hash, hash[:]) &&
			time.Now().Before(token.Expiry) {
			return token, nil
		}
	}

	return nil, nil
}

