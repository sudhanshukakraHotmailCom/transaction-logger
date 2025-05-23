package auth_test

import (
	"testing"
	"time"


	"transaction-logger/internal/auth"
	"transaction-logger/internal/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAndValidateJWT(t *testing.T) {
	// Setup test config
	cfg := &config.Config{
		JWTSecret: "test-secret-key-123",
	}
	auth.Init(cfg)

	userID := "test-user-123"
	email := "test@example.com"

	// Test token generation
	tokenString, err := auth.GenerateJWT(userID, email)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Test token validation
	t.Run("valid token", func(t *testing.T) {
		token, err := jwt.ParseWithClaims(tokenString, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})
		assert.NoError(t, err)
		
		if claims, ok := token.Claims.(*auth.Claims); ok && token.Valid {
			assert.Equal(t, userID, claims.UserID)
			assert.Equal(t, email, claims.Email)
		} else {
			t.Fatal("Invalid token claims")
		}
	})

	// Test expired token
	t.Run("expired token", func(t *testing.T) {
		// Create an expired token
		expiredToken, err := generateTestTokenWithExpiration(userID, email, -time.Hour, cfg.JWTSecret)
		assert.NoError(t, err)

		_, err = jwt.ParseWithClaims(expiredToken, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token is expired")
	})

	// Test invalid signature
	t.Run("invalid signature", func(t *testing.T) {
		// Create a token with different secret
		differentToken, err := generateTestTokenWithExpiration(userID, email, time.Hour, "different-secret-key")
		assert.NoError(t, err)

		_, err = jwt.ParseWithClaims(differentToken, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})
		assert.Error(t, err)
	})
}

// generateTestTokenWithExpiration is a test helper to generate a token with custom expiration
func generateTestTokenWithExpiration(userID, email string, expiration time.Duration, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &auth.Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
		},
	})

	return token.SignedString([]byte(secret))
}
