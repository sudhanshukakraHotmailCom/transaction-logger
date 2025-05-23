package testutils

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"transaction-logger/internal/database"
)

// SetupTestDB creates a new test database connection and runs migrations
func SetupTestDB(t *testing.T) *database.Database {
	t.Helper()
	
	testDB := os.Getenv("TEST_DB")
	if testDB == "" {
		t.Skip("Skipping integration test: TEST_DB environment variable not set")
	}

	db, err := sql.Open("postgres", testDB)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Create a new database instance
	dbInstance := &database.Database{DB: db}

	// Run migrations
	if err := dbInstance.InitSchema(); err != nil {
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

	return dbInstance
}

// CreateTestUser creates a test user and returns the user ID
func CreateTestUser(t *testing.T, db *database.Database, email, password string) string {
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
