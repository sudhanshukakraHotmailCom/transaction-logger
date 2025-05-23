package models_test

import (
	"testing"

	"transaction-logger/internal/models"
	"transaction-logger/internal/testutils"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	db := testutils.SetupTestDB(t)

	tests := []struct {
		name     string
		email    string
		password string
		wantErr  bool
	}{
		{"valid user", "test@example.com", "password123", false},
		{"duplicate email", "test@example.com", "password123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := models.CreateUser(db.DB, tt.email, tt.password)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, user.ID)
			assert.Equal(t, tt.email, user.Email)
			assert.NotEqual(t, tt.password, user.Password) // Password should be hashed
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	db := testutils.SetupTestDB(t)

	// Create a test user
	email := "test@example.com"
	password := "password123"
	_, err := models.CreateUser(db.DB, email, password)
	assert.NoError(t, err)

	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"existing user", email, false},
		{"non-existent user", "nonexistent@example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := models.GetUserByEmail(db.DB, tt.email)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.email, user.Email)
		})
	}
}
