# Test Setup and Execution Guide

## Prerequisites

- Go 1.21 or later
- PostgreSQL 13+ (for integration tests)
- Docker and Docker Compose (optional, for containerized testing)

## Test Structure

The test suite is organized as follows:

```
test/
└── integration/         # Integration tests
    ├── auth/           # Authentication tests
    └── models/         # Database model tests
        ├── transaction_test.go  # Transaction model tests including count verification
        └── user_test.go         # User model tests
```

### Key Test Cases

#### Transaction Tests
- Basic CRUD operations
- Transaction validation (amount, currency, etc.)
- **Transaction counting** - Verifies accurate counting of transactions after additions
- **Pagination** - Tests pagination functionality with various page sizes and edge cases
- User isolation - Ensures users can only access their own transactions
- Error handling for invalid inputs

#### Pagination Testing

Pagination is tested with the following scenarios:

1. **Default Pagination**
   - Verifies default page size (20 items)
   - Checks that first page is returned by default

2. **Custom Page Sizes**
   - Tests with different page sizes (within limits)
   - Verifies maximum page size limit (100 items)

3. **Page Navigation**
   - Tests navigation through multiple pages
   - Verifies `has_more` flag behavior
   - Tests last page handling

4. **Edge Cases**
   - Page number less than 1
   - Page number greater than total pages
   - Empty result sets
   - Single page results

Example test command for pagination tests:
```bash
# Run pagination tests
go test -v -run TestTransactionPagination ./test/integration/models/...
```

#### User Tests
- User registration and authentication
- Password hashing and verification
- Duplicate email prevention
- User data validation

## Running Tests

### Unit Tests

Run all unit tests:

```bash
go test -v ./...
```

### Integration Tests

1. **Set up test database**
   - Create a test database in PostgreSQL:
     ```sql
     CREATE DATABASE transaction_logger_test;
     ```

2. **Set environment variables**
   ```bash
   export TEST_DB=postgres://username:password@localhost:5432/transaction_logger_test?sslmode=disable
   ```

3. **Run integration tests**
   ```bash
   # Run all integration tests
   go test -v ./test/integration/...
   
   # Run specific test package
   go test -v ./test/integration/auth/...
   
   # Run specific test function
   go test -v -run TestCreateTransaction ./test/integration/models/...
   
   # Run transaction count test
   go test -v -run TestTransactionCountAfterAddition ./test/integration/models/...
   ```

### Running Tests with Docker

1. Start the test database:
   ```bash
   docker-compose -f docker-compose.test.yml up -d
   ```

2. Run tests:
   ```bash
   docker-compose -f docker-compose.test.yml run --rm app go test -v ./...
   ```

## Testing Pagination

### Unit Tests
Unit tests for pagination logic are located in `internal/utils/pagination_test.go`. These tests verify:
- Default page and page size values
- Offset calculation
- Parameter validation

### Integration Tests
Integration tests verify the pagination behavior with the actual database. These tests are in `test/integration/models/transaction_test.go` and test:
- Database query construction with LIMIT and OFFSET
- Pagination metadata in API responses
- Transaction isolation between users

### Testing with cURL
You can test pagination directly using cURL:

```bash
# Get first page with 10 items per page
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" "http://localhost:8080/api/transactions?page=1&page_size=10"

# Get second page
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" "http://localhost:8080/api/transactions?page=2&page_size=10"
```

## Test Coverage

Generate coverage report:

```bash
# Generate coverage profile
mkdir -p coverage
go test -coverprofile=coverage/coverage.out ./...

# Generate HTML report
go tool cover -html=coverage/coverage.out -o coverage/coverage.html

# View coverage in terminal
go tool cover -func=coverage/coverage.out
```

## Test Data Management

### Test Database Setup

The test database is automatically set up and torn down using the `testutils` package. Each test runs in its own transaction that is rolled back after the test completes.

### Test Fixtures

Common test data is defined in the test files. The following helper functions are available:

- `testutils.CreateTestUser()` - Create a test user
- `testutils.CreateTestTransaction()` - Create a test transaction
- `testutils.CleanupTestData()` - Clean up test data

### Test Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `TEST_DB` | Test database connection string | Required |
| `TEST_LOG_LEVEL` | Log level for tests | `error` |
| `TEST_TIMEOUT` | Test timeout duration | `30s` |

## Debugging Tests

### Verbose Output

Add `-v` flag for verbose output:
```bash
go test -v ./...
```

### Debug with Delve

1. Install Delve:
   ```bash
   go install github.com/go-delve/delve/cmd/dlv@latest
   ```

2. Debug tests:
   ```bash
   dlv test -- -test.run TestFunctionName
   ```

## Common Issues

1. **Database Connection Issues**
   - Ensure PostgreSQL is running
   - Verify database credentials in the connection string
   - Check if the test database exists and is accessible

2. **Test Failures**
   - Run with `-v` flag for more detailed output
   - Check for database constraint violations
   - Verify test data setup and cleanup

3. **Race Conditions**
   Run tests with the race detector:
   ```bash
   go test -race ./...
   ```

## Performance Testing

Run benchmarks:
```bash
# Run all benchmarks
go test -bench=. ./...

# Run specific benchmark
go test -bench=BenchmarkTransactionCreate ./...

# Run with profiling
# CPU profile
go test -bench=. -cpuprofile=cpu.out ./...
# Memory profile
go test -bench=. -memprofile=mem.out ./...
```

## Continuous Integration

The test suite is configured to run on pull requests and pushes to the main branch. The CI pipeline includes:
- Unit tests
- Integration tests
- Race detection
- Code coverage reporting
- Linting
- Security scanning
