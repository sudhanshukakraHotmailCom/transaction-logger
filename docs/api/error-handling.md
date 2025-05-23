# Error Handling

This document outlines the standard error responses and status codes returned by the Transaction Service API.

## Standard Error Response Format

All error responses follow this format:

```json
{
  "error": "Error Type",
  "message": "Human-readable error message",
  "details": {
    "field_name": "Specific error message"
  }
}
```

## Common HTTP Status Codes

### 400 Bad Request
- **When**: The request contains invalid data or missing required fields
- **Example**:
  ```json
  {
    "error": "Bad Request",
    "message": "Invalid request data",
    "details": {
      "email": "must be a valid email address"
    }
  }
  ```

### 401 Unauthorized
- **When**: Missing or invalid authentication token
- **Example**:
  ```json
  {
    "error": "Unauthorized",
    "message": "Authentication required"
  }
  ```

### 403 Forbidden
- **When**: Authenticated user doesn't have permission
- **Example**:
  ```json
  {
    "error": "Forbidden",
    "message": "You don't have permission to access this resource"
  }
  ```

### 404 Not Found
- **When**: Requested resource doesn't exist
- **Example**:
  ```json
  {
    "error": "Not Found",
    "message": "Transaction not found"
  }
  ```

### 409 Conflict
- **When**: Resource conflict (e.g., duplicate email)
- **Example**:
  ```json
  {
    "error": "Conflict",
    "message": "Email already registered"
  }
  ```

### 422 Unprocessable Entity
- **When**: Request is well-formed but contains semantic errors
- **Example**:
  ```json
  {
    "error": "Unprocessable Entity",
    "message": "Insufficient funds"
  }
  ```

### 500 Internal Server Error
- **When**: Unexpected server error
- **Example**:
  ```json
  {
    "error": "Internal Server Error",
    "message": "Something went wrong"
  }
  ```

## Validation Errors

When request validation fails, the API returns a 400 Bad Request with details about the validation errors:

```json
{
  "error": "Bad Request",
  "message": "Validation failed",
  "details": {
    "email": "must be a valid email address",
    "password": "must be at least 8 characters long"
  }
}
```

## Rate Limiting

- **Status Code**: 429 Too Many Requests
- **Headers**:
  - `Retry-After`: Number of seconds to wait before making a new request
- **Example**:
  ```json
  {
    "error": "Too Many Requests",
    "message": "Rate limit exceeded. Please try again in 60 seconds"
  }
  ```
