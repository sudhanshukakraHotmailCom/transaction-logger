# Authentication API

## Register a New User

### Endpoint
```
POST /register
```

### Description
Creates a new user account and returns a JWT token for authentication.

### Request
```http
POST /register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

### Request Body
| Field    | Type   | Required | Description          |
|----------|--------|----------|----------------------|
| email    | string | Yes      | User's email address |
| password | string | Yes      | User's password      |

### Response
#### Success (200 OK)
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### Error (400 Bad Request)
```json
{
  "error": "Bad Request",
  "message": "Email already registered"
}
```

## User Login

### Endpoint
```
POST /login
```

### Description
Authenticates a user and returns a JWT token.

### Request
```http
POST /login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

### Request Body
| Field    | Type   | Required | Description          |
|----------|--------|----------|----------------------|
| email    | string | Yes      | User's email address |
| password | string | Yes      | User's password      |

### Response
#### Success (200 OK)
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### Error (401 Unauthorized)
```json
{
  "error": "Unauthorized",
  "message": "Invalid email or password"
}
```
