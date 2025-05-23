package handlers

import (
	"database/sql"
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"transaction-logger/internal/models"
)

type TransactionHandler struct {
	db *sql.DB
}

func NewTransactionHandler(db *sql.DB) *TransactionHandler {
	return &TransactionHandler{db: db}
}

func (h *TransactionHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by AuthMiddleware)
	userID := r.Context().Value("userID").(string)

	rows, err := h.db.Query(
		`SELECT id, timestamp, sender_account, receiver_account, 
		  amount, currency, transaction_type, status, user_id 
		  FROM transactions WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(
			&t.ID,
			&t.Timestamp,
			&t.SenderAccount,
			&t.ReceiverAccount,
			&t.Amount,
			&t.Currency,
			&t.TransactionType,
			&t.Status,
			&t.UserID,
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		transactions = append(transactions, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}



// CreateTransaction handles the creation of a single transaction
func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by AuthMiddleware)
	userID := r.Context().Value("userID").(string)

	var req models.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set the user ID from the authenticated user
	req.UserID = userID

	// Validate request
	if req.SenderAccount == "" || req.ReceiverAccount == "" || req.Amount <= 0 {
		http.Error(w, "Invalid transaction data", http.StatusBadRequest)
		return
	}

	tx := models.Transaction{
		ID:              generateID(),
		Timestamp:       time.Now(),
		SenderAccount:   req.SenderAccount,
		ReceiverAccount: req.ReceiverAccount,
		Amount:          req.Amount,
		Currency:        req.Currency,
		TransactionType: req.TransactionType,
		Status:          "Completed",
		UserID:          userID,
	}

	_, err := h.db.Exec(
		`INSERT INTO transactions 
		(id, timestamp, sender_account, receiver_account, amount, currency, transaction_type, status, user_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		tx.ID, tx.Timestamp, tx.SenderAccount, tx.ReceiverAccount, tx.Amount, tx.Currency, tx.TransactionType, tx.Status, tx.UserID,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tx)
}

// GenerateSampleTransactions generates sample transactions for testing
func (h *TransactionHandler) GenerateSampleTransactions(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by AuthMiddleware)
	userID := r.Context().Value("userID").(string)

	tx, err := h.db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Prepare the statement with user_id parameter
	stmt, err := tx.Prepare(`
		INSERT INTO transactions (
			id, timestamp, sender_account, receiver_account, 
			amount, currency, transaction_type, status, user_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Generate 100 sample transactions
	for i := 0; i < 100; i++ {
		timestamp := time.Now().Add(-time.Duration(rand.Intn(365)) * 24 * time.Hour)
		sender := generateAccountNumber()
		receiver := generateAccountNumber()
		amount := float64(rand.Intn(10000)) + rand.Float64()
		currency := []string{"USD", "EUR", "GBP"}[rand.Intn(3)]
		txType := []string{"Transfer", "Deposit", "Withdrawal"}[rand.Intn(3)]

		_, err = stmt.Exec(
			generateID(),
			timestamp,
			sender,
			receiver,
			amount,
			currency,
			txType,
			"Completed",
			userID, // Include the user ID in the transaction
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Successfully generated 100 transactions"))
}

func generateID() string {
	rand.Seed(time.Now().UnixNano())
	return "TXN" + time.Now().Format("20060102150405") + strconv.Itoa(rand.Intn(1000))
}

// generateAccountNumber generates a random 12-digit account number
func generateAccountNumber() string {
	const digits = "0123456789"
	b := make([]byte, 12) // 12-digit account number
	for i := range b {
		b[i] = digits[rand.Intn(len(digits))]
	}
	return string(b)
}
