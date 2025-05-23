package models_test

import (
	"testing"
	"time"

	"transaction-logger/internal/models"
	"transaction-logger/internal/testutils"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransaction(t *testing.T) {
	db := testutils.SetupTestDB(t)

	// Create a test user
	userID := testutils.CreateTestUser(t, db, "test@example.com", "password123")


	tests := []struct {
		name          string
		transaction   models.CreateTransactionRequest
		userID        string
		wantErr       bool
		expectedError string
	}{
		{
			name: "valid transaction",
			transaction: models.CreateTransactionRequest{
				SenderAccount:   "ACC12345678",
				ReceiverAccount: "ACC87654321",
				Amount:          100.50,
				Currency:        "USD",
				TransactionType: "Transfer",
			},
			userID:  userID,
			wantErr: false,
		},
		{
			name: "invalid amount",
			transaction: models.CreateTransactionRequest{
				SenderAccount:   "ACC12345678",
				ReceiverAccount: "ACC87654321",
				Amount:          0, // Invalid amount
				Currency:        "USD",
				TransactionType: "Transfer",
			},
			userID:        userID,
			wantErr:       true,
			expectedError: "amount must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := models.CreateTransaction(db.DB, tt.transaction, tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedError != "" {
					assert.Contains(t, err.Error(), tt.expectedError)
				}
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, tx.ID)
			assert.Equal(t, tt.transaction.SenderAccount, tx.SenderAccount)
			assert.Equal(t, tt.transaction.ReceiverAccount, tx.ReceiverAccount)
			assert.Equal(t, tt.transaction.Amount, tx.Amount)
			assert.Equal(t, tt.transaction.Currency, tx.Currency)
			assert.Equal(t, tt.transaction.TransactionType, tx.TransactionType)
			assert.Equal(t, "Completed", tx.Status)
			assert.Equal(t, tt.userID, tx.UserID)
		})
	}
}

func TestGetTransactions(t *testing.T) {
	db := testutils.SetupTestDB(t)

	// Create a test user
	userID := testutils.CreateTestUser(t, db, "test@example.com", "password123")

	// Create some test transactions
	transactions := []models.CreateTransactionRequest{
		{
			SenderAccount:   "ACC11111111",
			ReceiverAccount: "ACC22222222",
			Amount:          100.00,
			Currency:        "USD",
			TransactionType: "Transfer",
		},
		{
			SenderAccount:   "ACC33333333",
			ReceiverAccount: "ACC44444444",
			Amount:          200.00,
			Currency:        "EUR",
			TransactionType: "Deposit",
		},
	}

	// Insert test transactions
	for _, tx := range transactions {
		_, err := models.CreateTransaction(db.DB, tx, userID)
		assert.NoError(t, err)
	}

	// Test getting transactions
	t.Run("get all transactions for user", func(t *testing.T) {
		txs, err := models.GetTransactions(db.DB, userID)
		assert.NoError(t, err)
		assert.Len(t, txs, len(transactions))

		// Verify the transactions match what we created
		for i, tx := range txs {
			assert.Equal(t, transactions[i].SenderAccount, tx.SenderAccount)
			assert.Equal(t, transactions[i].ReceiverAccount, tx.ReceiverAccount)
			assert.Equal(t, transactions[i].Amount, tx.Amount)
			assert.Equal(t, transactions[i].Currency, tx.Currency)
			assert.Equal(t, transactions[i].TransactionType, tx.TransactionType)
			assert.Equal(t, userID, tx.UserID)
		}
	})

	// Test getting transactions for non-existent user
	t.Run("get transactions for non-existent user", func(t *testing.T) {
		txs, err := models.GetTransactions(db.DB, "non-existent-user-id")
		assert.NoError(t, err)
		assert.Empty(t, txs)
	})
}
