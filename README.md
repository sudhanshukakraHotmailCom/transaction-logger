# Transaction Logger Microservice

A Go-based microservice for logging and managing financial transactions with a PostgreSQL backend, featuring JWT-based authentication.

## Quick Start

### Prerequisites

- Docker and Docker Compose
- (Optional) Go 1.21+ for local development

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

### Manual Setup

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Set up environment variables:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. Start the database:
   ```bash
   docker-compose up -d db
   ```

4. Run migrations:
   ```bash
   go run cmd/server/main.go migrate
   ```

5. Start the server:
   ```bash
   go run cmd/server/main.go serve
   ```

## API Usage

### Authentication

All API endpoints except `/register` and `/login` require a valid JWT token:

```
Authorization: Bearer <your-jwt-token>
```

### Available Endpoints

#### Authentication
- `POST /register` - Register a new user
- `POST /login` - Authenticate and get JWT token

#### Features

- User authentication with JWT
- Create and retrieve transactions with pagination support
- Transaction validation
- Sample data generation for testing
- RESTful API endpoints

#### Transactions
- `POST /transactions` - Create a new transaction
- `GET /transactions` - List all transactions for the authenticated user

## Pagination

The API supports pagination for transaction listings with the following query parameters:

- `page` - The page number to retrieve (default: 1)
- `page_size` - Number of items per page (default: 20, max: 100)

Example request:
```http
GET /api/transactions?page=2&page_size=10
```

Example response:
```json
{
  "data": [
    {
      "id": "...",
      "timestamp": "...",
      "sender_account": "...",
      "receiver_account": "...",
      "amount": 100.5,
      "currency": "USD",
      "transaction_type": "Transfer",
      "status": "Completed"
    }
  ],
  "pagination": {
    "total": 42,
    "count": 10,
    "per_page": 10,
    "current_page": 2,
    "total_pages": 5,
    "has_more": true
  }
}
```

### Example: Create a Transaction

```bash
curl -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "sender_account": "ACC123456",
    "receiver_account": "ACC789012",
    "amount": 100.50,
    "currency": "USD",
    "transaction_type": "Transfer"
  }'
```

## Documentation

- [Detailed Documentation](DETAILED_DOCUMENTATION.md) - Project structure, components, and API reference
- [Test Setup](TEST_SETUP.md) - Testing instructions and guidelines

## Configuration

### Environment Variables

#### Database
- `POSTGRES_HOST`: Database host (default: localhost)
- `POSTGRES_PORT`: Database port (default: 5432)
- `POSTGRES_USER`: Database user (default: postgres)
- `POSTGRES_PASSWORD`: Database password (default: postgres)
- `POSTGRES_DB`: Database name (default: transaction_logger)

#### Server
- `PORT`: HTTP server port (default: 8080)
- `JWT_SECRET`: Secret key for JWT token generation (required in production)

## License

MIT License - See [LICENSE](LICENSE) for details.