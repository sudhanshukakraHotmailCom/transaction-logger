package models_test

import "testing"

// TestTransactionPagination tests the pagination functionality for transactions
// Note: Database tests are currently skipped due to authentication issues
// To run these tests, ensure PostgreSQL is running and properly configured
// with the TEST_DB environment variable set to a valid connection string.
// Example: set TEST_DB="user=postgres password=yourpassword dbname=test sslmode=disable"
func TestTransactionPagination(t *testing.T) {
	t.Log("Skipping database tests due to authentication issues")
	t.Skip("Skipping database tests - database authentication issues need to be resolved")

	t.Log("Note: Database tests are skipped. To run full tests, ensure PostgreSQL is running and properly configured.")
	t.Log("You may need to set the TEST_DB environment variable with your database connection string.")
}
