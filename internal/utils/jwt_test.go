package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateJWT(t *testing.T) {
	tests := []struct {
		name     string
		userID   uint
		username string
		email    string
		secret   string
		wantErr  bool
	}{
		{
			name:     "Valid JWT generation",
			userID:   1,
			username: "testuser",
			email:    "test@example.com",
			secret:   "test-secret",
			wantErr:  false,
		},
		{
			name:     "Empty secret",
			userID:   1,
			username: "testuser",
			email:    "test@example.com",
			secret:   "",
			wantErr:  false, // JWT library allows empty secret
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateJWT(tt.userID, tt.username, tt.email, tt.secret)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				
				// Verify token can be parsed
				parsedToken, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
					return []byte(tt.secret), nil
				})
				
				assert.NoError(t, err)
				assert.True(t, parsedToken.Valid)
				
				claims, ok := parsedToken.Claims.(*JWTClaims)
				assert.True(t, ok)
				assert.Equal(t, tt.userID, claims.UserID)
				assert.Equal(t, tt.username, claims.Username)
				assert.Equal(t, tt.email, claims.Email)
			}
		})
	}
}

func TestGenerateTokenPair(t *testing.T) {
	tests := []struct {
		name     string
		userID   uint
		username string
		email    string
		secret   string
		wantErr  bool
	}{
		{
			name:     "Valid token pair generation",
			userID:   1,
			username: "testuser",
			email:    "test@example.com",
			secret:   "test-secret",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenPair, err := GenerateTokenPair(tt.userID, tt.username, tt.email, tt.secret)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, tokenPair)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tokenPair)
				assert.NotEmpty(t, tokenPair.AccessToken)
				assert.NotEmpty(t, tokenPair.RefreshToken)
				assert.Greater(t, tokenPair.ExpiresIn, int64(0))
				
				// Verify both tokens can be validated
				accessClaims, err := ValidateJWT(tokenPair.AccessToken, tt.secret)
				assert.NoError(t, err)
				assert.Equal(t, tt.userID, accessClaims.UserID)
				
				refreshClaims, err := ValidateJWT(tokenPair.RefreshToken, tt.secret)
				assert.NoError(t, err)
				assert.Equal(t, tt.userID, refreshClaims.UserID)
				
				// Verify refresh token has longer expiry than access token
				assert.True(t, refreshClaims.ExpiresAt.After(*accessClaims.ExpiresAt))
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	secret := "test-secret"
	userID := uint(1)
	username := "testuser"
	email := "test@example.com"
	
	// Generate a valid token
	validToken, err := GenerateJWT(userID, username, email, secret)
	require.NoError(t, err)
	
	// Generate an expired token
	expiredClaims := JWTClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)), // Expired 1 hour ago
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-time.Hour * 2)),
		},
	}
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredTokenString, err := expiredToken.SignedString([]byte(secret))
	require.NoError(t, err)

	tests := []struct {
		name        string
		token       string
		secret      string
		wantErr     bool
		expectedID  uint
	}{
		{
			name:       "Valid token",
			token:      validToken,
			secret:     secret,
			wantErr:    false,
			expectedID: userID,
		},
		{
			name:    "Invalid token format",
			token:   "invalid-token",
			secret:  secret,
			wantErr: true,
		},
		{
			name:    "Wrong secret",
			token:   validToken,
			secret:  "wrong-secret",
			wantErr: true,
		},
		{
			name:    "Expired token",
			token:   expiredTokenString,
			secret:  secret,
			wantErr: true,
		},
		{
			name:    "Empty token",
			token:   "",
			secret:  secret,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateJWT(tt.token, tt.secret)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, tt.expectedID, claims.UserID)
				assert.Equal(t, username, claims.Username)
				assert.Equal(t, email, claims.Email)
			}
		})
	}
}

func TestExtractUserID(t *testing.T) {
	secret := "test-secret"
	userID := uint(123)
	username := "testuser"
	email := "test@example.com"
	
	// Generate a valid token
	validToken, err := GenerateJWT(userID, username, email, secret)
	require.NoError(t, err)

	tests := []struct {
		name       string
		token      string
		secret     string
		wantErr    bool
		expectedID uint
	}{
		{
			name:       "Valid token",
			token:      validToken,
			secret:     secret,
			wantErr:    false,
			expectedID: userID,
		},
		{
			name:    "Invalid token",
			token:   "invalid-token",
			secret:  secret,
			wantErr: true,
		},
		{
			name:    "Wrong secret",
			token:   validToken,
			secret:  "wrong-secret",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extractedID, err := ExtractUserID(tt.token, tt.secret)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, uint(0), extractedID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, extractedID)
			}
		})
	}
}

func TestJWTClaims(t *testing.T) {
	claims := &JWTClaims{
		UserID:   1,
		Username: "testuser",
		Email:    "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	assert.Equal(t, uint(1), claims.UserID)
	assert.Equal(t, "testuser", claims.Username)
	assert.Equal(t, "test@example.com", claims.Email)
	assert.NotNil(t, claims.ExpiresAt)
	assert.NotNil(t, claims.IssuedAt)
}

func TestTokenPair(t *testing.T) {
	tokenPair := &TokenPair{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresIn:    3600,
	}

	assert.Equal(t, "access-token", tokenPair.AccessToken)
	assert.Equal(t, "refresh-token", tokenPair.RefreshToken)
	assert.Equal(t, int64(3600), tokenPair.ExpiresIn)
}