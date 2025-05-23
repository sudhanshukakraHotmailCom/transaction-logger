# Transaction Logger Microservice

A Go-based microservice for logging and managing financial transactions with a PostgreSQL backend, featuring JWT-based authentication.

## Features

- User registration and authentication with JWT
- Secure transaction management with user isolation
- Create and store transaction records
- Generate sample transaction data
- RESTful API for transaction management
- Containerized with Docker
- PostgreSQL database integration

## Prerequisites

- Docker and Docker Compose
- Go 1.21 or later (only needed for local development without Docker)

## Getting Started

### Using Docker Compose (Recommended)

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd transaction-logger
   ```

2. Start the services:
   ```bash
   docker-compose up --build
   ```

   This will:
   - Build the Go application
   - Start a PostgreSQL database
   - Run the transaction logger service on port 8080

## Authentication

All API endpoints except `/register` and `/login` require a valid JWT token in the `Authorization` header:

```
Authorization: Bearer <your-jwt-token>
```

### API Endpoints

#### Authentication
- `POST /register` - Register a new user
- `POST /login` - Login and get a JWT token

#### Transactions (Requires Authentication)
- `GET /transactions` - Get all transactions for the authenticated user
- `POST /transactions` - Create a new transaction for the authenticated user
- `POST /transactions/generate` - Generate sample transactions for the authenticated user

#### Example: Register a New User

```http
POST /register HTTP/1.1
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

#### Example: Login

```http
POST /login HTTP/1.1
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

#### Example: Create a Transaction (Authenticated)

```http
POST /transactions HTTP/1.1
Content-Type: application/json
Authorization: Bearer <your-jwt-token>

{
  "sender_account": "ACC12345678",
  "receiver_account": "ACC87654321",
  "amount": 100.50,
  "currency": "USD",
  "transaction_type": "Transfer"
}
```

Note: The `status` field is automatically set to "Completed" by the server.

## Development

### Running Tests

```bash
go test -v ./...
```

### Environment Variables

### Database
- `POSTGRES_HOST`: Database host (default: localhost)
- `POSTGRES_PORT`: Database port (default: 5432)
- `POSTGRES_USER`: Database user (default: postgres)
- `POSTGRES_PASSWORD`: Database password (default: postgres)
- `POSTGRES_DB`: Database name (default: transaction_logger)

### Server
- `PORT`: HTTP server port (default: 8080)
- `JWT_SECRET`: Secret key for JWT token generation (required in production)

## Running Migrations

Database migrations are automatically applied when the application starts. For manual migrations:

```bash
# Apply migrations
cat migrations/*.up.sql | psql $DATABASE_URL

# Rollback last migration (if needed)
cat migrations/002_add_user_id_to_transactions.down.sql | psql $DATABASE_URL
```

## License

MIT