package testutils

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

// TestDB wraps a database connection for testing
type TestDB struct {
	*sql.DB
}

// SetupTestDB creates a new test database connection and runs migrations
func SetupTestDB(t *testing.T) *TestDB {
	t.Helper()
	
	testDB := os.Getenv("TEST_DB")
	if testDB == "" {
		t.Skip("Skipping integration test: TEST_DB environment variable not set")
	}

	db, err := sql.Open("postgres", testDB)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Initialize database schema
	if err := initDBSchema(db); err != nil {
		t.Fatalf("Failed to initialize test database schema: %v", err)
	}

	t.Cleanup(func() {
		// Clean up test data
		_, err := db.Exec(`
			DROP TABLE IF EXISTS transactions CASCADE;
			DROP TABLE IF EXISTS users CASCADE;
		`)
		if err != nil {
			t.Logf("Failed to clean up test database: %v", err)
		}
		db.Close()
	})

	return &TestDB{db}
}

// initDBSchema initializes the database schema
func initDBSchema(db *sql.DB) error {
	// Create users table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL,
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL
		);
	`)
	if err != nil {
		return err
	}

	// Create transactions table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS transactions (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			sender_account TEXT NOT NULL,
			receiver_account TEXT NOT NULL,
			amount DECIMAL(15, 2) NOT NULL,
			currency TEXT NOT NULL,
			transaction_type TEXT NOT NULL,
			status TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL,
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL
		);
	`)

	return err
}

// CreateTestUser creates a test user and returns the user ID
func CreateTestUser(t *testing.T, db *TestDB, email, password string) string {
	t.Helper()

	var userID string
	err := db.DB.QueryRow(
		`INSERT INTO users (id, email, password, created_at, updated_at)
		 VALUES ($1, $2, $3, NOW(), NOW())
		 RETURNING id`,
		"test-user-123", email, password,
	).Scan(&userID)

	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	return userID
}
