# Transaction Logger Microservice - Detailed Documentation

## Project Overview

A Go-based microservice for logging and managing financial transactions with a PostgreSQL backend, featuring JWT-based authentication.

## Project Structure

The project follows a clean architecture pattern with the following directory structure:

```
.
├── cmd/
│   └── server/           # Application entry point
│       └── main.go       # Main application setup and server initialization
│
├── internal/
│   ├── auth/            # Authentication related code
│   │   └── auth.go      # JWT authentication middleware and utilities
│   │
│   ├── config/          # Configuration management
│   │   └── config.go    # Environment variables and configuration
│   │
│   ├── database/        # Database connection and migrations
│   │   └── database.go  # Database setup and connection management
│   │
│   ├── handlers/        # HTTP request handlers
│   │   ├── auth.go      # Authentication handlers (register, login)
│   │   └── transaction.go # Transaction management handlers
│   │
│   ├── models/          # Data models and database operations
│   │   ├── transaction.go # Transaction model and database operations
│   │   └── user.go       # User model and database operations
│   │
│   └── testutils/       # Test utilities
│       └── testutils.go  # Test database setup and helper functions
│
├── test/
│   └── integration/    # Integration tests
│       ├── auth/        # Authentication tests
│       └── models/      # Model tests
```

## Component Documentation

### 1. Pagination Implementation (`internal/utils/pagination.go`)

The pagination system provides a consistent way to handle large sets of data across the application.

#### Key Components:

1. **Pagination Struct**
   ```go
   type Pagination struct {
       Page     int
       PageSize int
   }
   ```
   - `Page`: Current page number (1-based)
   - `PageSize`: Number of items per page

2. **NewPagination**
   - Creates a new Pagination instance from HTTP request parameters
   - Sets defaults (page=1, page_size=20)
   - Enforces maximum page size (100 items)

3. **Offset**
   - Calculates the database offset for SQL queries
   - Implements: `(page - 1) * pageSize`

#### Usage in Handlers

```go
// Example usage in a handler
func (h *Handler) ListItems(c *gin.Context) {
    // Get pagination from query params
    pagination := utils.NewPagination(c)
    
    // Query database with pagination
    items, total, err := h.repo.GetItems(
        pagination.Offset(),
        pagination.PageSize,
    )
    
    // Calculate total pages
    totalPages := total / pagination.PageSize
    if total%pagination.PageSize > 0 {
        totalPages++
    }
    
    // Return paginated response
    c.JSON(200, gin.H{
        "data": items,
        "pagination": {
            "total": total,
            "count": len(items),
            "per_page": pagination.PageSize,
            "current_page": pagination.Page,
            "total_pages": totalPages,
            "has_more": pagination.Page < totalPages,
        },
    })
}
```

### 2. Authentication (`internal/auth`)

Handles JWT-based authentication:
- Token generation and validation
- Authentication middleware
- Password hashing and verification

### 2. Configuration (`internal/config`)

Manages application configuration:
- Environment variables
- Database connection settings
- JWT secret key
- Server configuration

### 3. Database (`internal/database`)

Database connection and migration management:
- PostgreSQL connection setup
- Database migration utilities
- Connection pooling configuration

### 4. Handlers (`internal/handlers`)

HTTP request handlers:
- `auth.go`: User registration and login endpoints
- `transaction.go`: Transaction CRUD operations
  - Create new transactions
  - List user transactions
  - Get transaction details

### 5. Models (`internal/models`)

Data models and database operations:
- `user.go`: User model and database operations
  - User creation and validation
  - Password hashing
  - User lookup by email/ID
- `transaction.go`: Transaction model and operations
  - Transaction creation and validation
  - Transaction queries
  - Business logic for transaction processing

### 6. Test Utilities (`internal/testutils`)

Testing helpers:
- Test database setup and teardown
- Test user creation
- Test transaction creation
- Database cleanup utilities

## API Documentation

### Authentication

#### Register a New User
```
POST /register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

#### Login
```
POST /login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

### Transactions (Requires Authentication)

#### Create a New Transaction
```
POST /transactions
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "sender_account": "ACC123456",
  "receiver_account": "ACC789012",
  "amount": 100.50,
  "currency": "USD",
  "transaction_type": "Transfer"
}
```

#### List User's Transactions
```
GET /transactions
Authorization: Bearer <jwt-token>
```

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);
```

### Transactions Table
```sql
CREATE TABLE transactions (
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
```

## Error Handling

The API returns appropriate HTTP status codes and JSON error responses in the following format:

```json
{
  "error": "Error message describing the issue"
}
```

## Security Considerations

- All passwords are hashed using bcrypt before storage
- JWT tokens are used for authentication
- Database connection uses SSL in production
- Sensitive configuration is stored in environment variables
- Input validation is performed on all API endpoints

## Pagination Best Practices

### Performance Considerations

1. **Indexing**
   - Ensure proper indexing on columns used in WHERE and ORDER BY clauses
   - Example index for transactions:
     ```sql
     CREATE INDEX idx_transactions_user_timestamp 
     ON transactions(user_id, timestamp DESC);
     ```

2. **Query Optimization**
   - Use `EXPLAIN ANALYZE` to verify query plans
   - Avoid OFFSET with large page numbers (consider keyset pagination for very large datasets)

3. **Caching**
   - Consider caching frequently accessed pages
   - Implement cache invalidation on data changes

### Security Considerations

1. **Input Validation**
   - Validate page and page_size parameters
   - Enforce maximum page size limits
   - Sanitize all database inputs

2. **Rate Limiting**
   - Implement rate limiting to prevent abuse
   - Consider different limits for authenticated vs unauthenticated users

## Performance Considerations

- Database connection pooling is configured for optimal performance
- Indexes are created on frequently queried columns
- Transactions are used for data consistency
- Response compression is enabled for large payloads

## Monitoring and Logging

- Structured logging is implemented using a logging middleware
- Request/response logging is enabled for all endpoints
- Error logging includes stack traces for debugging
- Log levels can be configured (DEBUG, INFO, WARN, ERROR)

## Future Enhancements

1. Add transaction filtering and pagination
2. Implement transaction export functionality
3. Add webhook support for transaction events
4. Implement rate limiting
5. Add OpenAPI/Swagger documentation
6. Add metrics and monitoring
7. Implement audit logging
8. Add support for transaction attachments
9. Implement two-factor authentication
10. Add support for transaction categories and tags
