# Transactions API

## Create a New Transaction

### Endpoint
```
POST /api/transactions
```

### Description
Creates a new transaction record.

### Authentication
- **Required**: Yes
- **Type**: Bearer Token

### Request
```http
POST /api/transactions
Content-Type: application/json
Authorization: Bearer YOUR_JWT_TOKEN

{
  "sender_account": "ACCOUNT123",
  "receiver_account": "ACCOUNT456",
  "amount": 150.75,
  "currency": "USD",
  "transaction_type": "transfer",
  "status": "pending"
}
```

### Request Body
| Field             | Type   | Required | Description                                  |
|-------------------|--------|----------|----------------------------------------------|
| sender_account    | string | Yes      | Sender's account number                      |
| receiver_account  | string | Yes      | Receiver's account number                    |
| amount            | number | Yes      | Transaction amount (must be positive)        |
| currency          | string | Yes      | 3-letter currency code (e.g., USD, EUR)     |
| transaction_type  | string | No       | Type of transaction (e.g., transfer, payment)|
| status            | string | No       | Initial status (default: "pending")          |

### Response
#### Success (201 Created)
```json
{
  "id": "txn_1234567890",
  "user_id": "user_123",
  "sender_account": "ACCOUNT123",
  "receiver_account": "ACCOUNT456",
  "amount": 150.75,
  "currency": "USD",
  "transaction_type": "transfer",
  "status": "completed",
  "timestamp": "2025-05-23T18:57:45Z"
}
```

#### Error (400 Bad Request)
```json
{
  "error": "Bad Request",
  "message": "Invalid transaction data",
  "details": {
    "amount": "must be greater than 0"
  }
}
```

## List Transactions

### Endpoint
```
GET /api/transactions
```

### Description
Retrieves a paginated list of transactions for the authenticated user.

### Authentication
- **Required**: Yes
- **Type**: Bearer Token

### Query Parameters
| Parameter | Type    | Required | Default | Description                     |
|-----------|---------|----------|---------|---------------------------------|
| page      | integer | No       | 1       | Page number (1-based)           |
| page_size | integer | No       | 20      | Number of items per page (max 100)|

### Request
```http
GET /api/transactions?page=1&page_size=10
Authorization: Bearer YOUR_JWT_TOKEN
```

### Response
#### Success (200 OK)
```json
{
  "data": [
    {
      "id": "txn_1234567890",
      "user_id": "user_123",
      "sender_account": "ACCOUNT123",
      "receiver_account": "ACCOUNT456",
      "amount": 150.75,
      "currency": "USD",
      "transaction_type": "transfer",
      "status": "completed",
      "timestamp": "2025-05-23T18:57:45Z"
    }
  ],
  "pagination": {
    "total": 42,
    "count": 10,
    "per_page": 10,
    "current_page": 1,
    "total_pages": 5,
    "has_more": true
  }
}
```

#### Error (401 Unauthorized)
```json
{
  "error": "Unauthorized",
  "message": "Authentication required"
}
```

## Get Transaction by ID

### Endpoint
```
GET /api/transactions/:id
```

### Description
Retrieves a specific transaction by its ID.

### Authentication
- **Required**: Yes
- **Type**: Bearer Token

### Path Parameters
| Parameter | Type   | Required | Description          |
|-----------|--------|----------|----------------------|
| id        | string | Yes      | Transaction ID       |

### Request
```http
GET /api/transactions/txn_1234567890
Authorization: Bearer YOUR_JWT_TOKEN
```

### Response
#### Success (200 OK)
```json
{
  "id": "txn_1234567890",
  "user_id": "user_123",
  "sender_account": "ACCOUNT123",
  "receiver_account": "ACCOUNT456",
  "amount": 150.75,
  "currency": "USD",
  "transaction_type": "transfer",
  "status": "completed",
  "timestamp": "2025-05-23T18:57:45Z",
  "created_at": "2025-05-23T18:57:45Z",
  "updated_at": "2025-05-23T18:57:45Z"
}
```

#### Error (404 Not Found)
```json
{
  "error": "Not Found",
  "message": "Transaction not found"
}
```
