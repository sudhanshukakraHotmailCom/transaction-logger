# Getting Started with Transaction Service API

This guide will help you get started with the Transaction Service API, covering authentication, making your first request, and handling responses.

## Base URL

All API endpoints are relative to the base URL:
```
http://localhost:8080
```

## Authentication

The API uses JWT (JSON Web Tokens) for authentication. Follow these steps to authenticate:

1. **Register a new user** (one-time setup)
   ```http
   POST /register
   Content-Type: application/json
   
   {
     "email": "user@example.com",
     "password": "securepassword123"
   }
   ```

2. **Login to get a token**
   ```http
   POST /login
   Content-Type: application/json
   
   {
     "email": "user@example.com",
     "password": "securepassword123"
   }
   ```
   
   Response:
   ```json
   {
     "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
   }
   ```

3. **Include the token in subsequent requests**
   ```
   Authorization: Bearer YOUR_JWT_TOKEN
   ```

## Making Your First Request

### Example: Create a Transaction

1. **Request**
   ```http
   POST /api/transactions
   Content-Type: application/json
   Authorization: Bearer YOUR_JWT_TOKEN
   
   {
     "sender_account": "ACCOUNT123",
     "receiver_account": "ACCOUNT456",
     "amount": 150.75,
     "currency": "USD",
     "transaction_type": "transfer"
   }
   ```

2. **Successful Response (201 Created)**
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

## Using cURL Examples

### Register a New User
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword123"
  }'
```

### Login and Save Token
```bash
# Login and save token to environment variable
export TOKEN=$(curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"securepassword123"}' | \
  jq -r '.token')

echo "Your token: $TOKEN"
```

### Create a Transaction
```bash
curl -X POST http://localhost:8080/api/transactions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "sender_account": "ACCOUNT123",
    "receiver_account": "ACCOUNT456",
    "amount": 150.75,
    "currency": "USD",
    "transaction_type": "transfer"
  }'
```

### List Transactions
```bash
curl "http://localhost:8080/api/transactions?page=1&page_size=5" \
  -H "Authorization: Bearer $TOKEN"
```

## SDKs and Client Libraries

### JavaScript/Node.js
```javascript
const axios = require('axios');

const api = axios.create({
  baseURL: 'http://localhost:8080',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${process.env.API_TOKEN}`
  }
});

// Example: Create a transaction
async function createTransaction(transactionData) {
  try {
    const response = await api.post('/api/transactions', transactionData);
    return response.data;
  } catch (error) {
    console.error('Error creating transaction:', error.response?.data || error.message);
    throw error;
  }
}
```

### Python
```python
import requests
import os

BASE_URL = 'http://localhost:8080'
HEADERS = {
    'Content-Type': 'application/json',
    'Authorization': f'Bearer {os.getenv("API_TOKEN")}'
}

def create_transaction(transaction_data):
    """Create a new transaction."""
    try:
        response = requests.post(
            f"{BASE_URL}/api/transactions",
            json=transaction_data,
            headers=HEADERS
        )
        response.raise_for_status()
        return response.json()
    except requests.exceptions.RequestException as e:
        print(f"Error creating transaction: {e}")
        raise
```

## Rate Limiting

The API is rate limited to prevent abuse. The current limits are:
- 100 requests per minute per IP address for unauthenticated endpoints
- 1000 requests per minute per user for authenticated endpoints

When rate limited, you'll receive a 429 status code with a `Retry-After` header indicating when you can retry.

## Best Practices

1. **Always use HTTPS** in production
2. **Never expose your API tokens** in client-side code
3. **Handle errors gracefully** - check status codes and error messages
4. **Implement retry logic** for transient failures
5. **Cache responses** when appropriate to reduce API calls
