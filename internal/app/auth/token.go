package auth

import (
	"sync"
	"time"
)

// TokenStore represents a simple in-memory token store
type TokenStore struct {
	tokens map[string]UserInfo
	mu     sync.RWMutex
}

// UserInfo represents user information stored in the token
type UserInfo struct {
	UserID   uint
	Username string
	Role     string
	Expiry   time.Time
}

// NewTokenStore creates a new token store
func NewTokenStore() *TokenStore {
	return &TokenStore{
		tokens: make(map[string]UserInfo),
	}
}

// StoreToken stores a token with user information
func (ts *TokenStore) StoreToken(token string, userInfo UserInfo) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.tokens[token] = userInfo
}

// GetToken retrieves token information
func (ts *TokenStore) GetToken(token string) (UserInfo, bool) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	userInfo, exists := ts.tokens[token]
	return userInfo, exists
}

// RemoveToken removes a token from the store
func (ts *TokenStore) RemoveToken(token string) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	delete(ts.tokens, token)
}

// CleanupExpiredTokens removes expired tokens
func (ts *TokenStore) CleanupExpiredTokens() {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	now := time.Now()
	for token, userInfo := range ts.tokens {
		if now.After(userInfo.Expiry) {
			delete(ts.tokens, token)
		}
	}
}

// Global token store instance
var TokenStoreInstance = NewTokenStore()

// GenerateToken generates a simple token with user information
func GenerateToken(userID uint, username, role string) string {
	token := "token_" + username + "_" + time.Now().Format("20060102150405")
	expiry := time.Now().Add(10 * 24 * time.Hour)
	
	TokenStoreInstance.StoreToken(token, UserInfo{
		UserID:   userID,
		Username: username,
		Role:     role,
		Expiry:   expiry,
	})
	
	return token
}

// ValidateToken validates a token and returns user information
func ValidateToken(token string) (UserInfo, bool) {
	userInfo, exists := TokenStoreInstance.GetToken(token)
	if !exists {
		return UserInfo{}, false
	}
	
	// Check if token is expired
	if time.Now().After(userInfo.Expiry) {
		TokenStoreInstance.RemoveToken(token)
		return UserInfo{}, false
	}
	
	return userInfo, true
}

// RemoveToken removes a token from the store
func RemoveToken(token string) {
	TokenStoreInstance.RemoveToken(token)
}