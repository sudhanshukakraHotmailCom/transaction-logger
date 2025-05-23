package auth_test

import (
	"testing"
	"time"

	"transaction-logger/internal/auth"
	"transaction-logger/internal/config"
	"transaction-logger/internal/testutils"

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
		claims, err := auth.ValidateJWT(tokenString)
		assert.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, email, claims.Email)
	})

	// Test expired token
	t.Run("expired token", func(t *testing.T) {
		// Create an expired token
		expiredToken, err := auth.GenerateTokenWithExpiration(userID, email, -time.Hour)
		assert.NoError(t, err)

		_, err = auth.ValidateJWT(expiredToken)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token has invalid claims: token is expired")
	})

	// Test invalid signature
	t.Run("invalid signature", func(t *testing.T) {
		// Create a token with different secret
		oldSecret := cfg.JWTSecret
		cfg.JWTSecret = "different-secret-key"
		auth.Init(cfg)
		differentToken, _ := auth.GenerateJWT(userID, email)

		// Restore original secret
		cfg.JWTSecret = oldSecret
		auth.Init(cfg)

		_, err = auth.ValidateJWT(differentToken)
		assert.Error(t, err)
	})
}

// Helper function to generate a token with custom expiration
generateTokenWithExpiration(userID, email string, expiration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
		},
	})

	return token.SignedString(jwtSecret)
}
