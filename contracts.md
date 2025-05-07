# Standard JSON Format:

## Errors:
```json
{
  "detail": {
    "status": status_code,
    "message": "A human-readable error message describing the problem."
  }
}
```

## Success:
```json
{
    "status": "success",
    "data": {}  // The data returned by the API.
}
```

## Pagination:
```json
{
    "status": "success",
    "data": {
        "page": page_number,
        "page_size": page_size,
        "total_items": total_items,
        "total_pages": total_pages,
        "filters": {},  // The filters used to retrieve the data (e.g., {"category": "expense", "date_gte": "2025-04-01"}).
        "items": []  // The list of items returned by the API.
    }
}
```

## Date and Time:
All dates and times will be represented in ISO 8601 format with UTC timezone: `YYYY-MM-DDTHH:MM:SSZ`.

## Authentication:
- Bearer token: The Bearer token should be included in the Authorization header of the HTTP request (e.g., Authorization: Bearer <your_access_token>).


# User Management Service
- Base path: /api/v1/

## User Model
- Attributes:
  - id: int
  - email: str
  - password: str
  - full_name: str

## Authentication Endpoints
- Base Path: /auth

### SignUp
- Endpoint: POST /signup
- Request Body:
  ```json
  {
    "email": "user@example.com",
    "password": "password123",
    "full_name": "John Doe",
  }
  ```
- Response:
  - Status Code: 201 Created
  ```json
  {
    "status": "success",
    "data": {
      "user_id": 123,
      "message": "User created successfully."
    }
  }
  ```


### LogIn
- Endpoint: POST /login
- Request Body:
  ```json
  {
    "email": "user@example.com",
    "password": "password123"
  }
  ```
  Response:
  - Status Code: 200 OK
  ```json
  {
    "access_token": "",  // The JWT token
    "expires_at": "2025-04-01T12:00:00Z",
    "refresh_token": ""  // The refresh token
  }
  ```

### Autheticate
- Endpoint: GET /authenticate?include=email,full_name
- Request Headers:
  ```json
  {
    "Authorization": "Bearer <your_access_token>"
  }
  ```
- Response:
  - Status Code: 200 OK
  ```json
  {
    "user_id": 1,
    "role": "user",
    "email": "user@example.com",
    "full_name": "John Doe"
  }
  ```

  Error Response:
  - Status Code: 401 Unauthorized
  ```json
  {
    "detail": {
      "status": 401,
      "message": "Invalid or expired token",
      }
  }
  ```
  - Headers:
  ```json
  {
    "WWW-Authenticate": "Bearer"
  }
  ```

#### Roles
The API will support the following roles:
- User

Future:
- Admin


## User Endpoints
- Base Path: /users

### Get Users
- Endpoint: GET /me **?include=field1,field2,...** (Optional query parameter to include specific fields in the response)
- Request Headers:
  ```json
  {
    "Authorization": "Bearer <your_access_token>"
  }
  ```
  Response:
  - Status Code: 200 OK
  ```json
  {
    "user_id": 1,
    "role": "user",
    "email": "user@example.com",
    "full_name": "John Doe"
  }
  ```

### Update User
- Endpoint: PATCH /me
- Request Headers:
  ```json
  {
    "Authorization": "Bearer <your_access_token>"
  }
  ```
- Request Body:
  ```json
  {
    "email": "user@example.com",
  }
  ``` OR
  ```json
  {
    "old_password": "password123",
    "new_password": "newpassword123"
  }
  ``` OR
  ```json
  {
    "full_name": "John Doe"
    // Other possible fields...
  }
  ```
  Response:
  - Status Code: 200 OK
  ```json
  {
    "status": "success",
    "data": {
      "user_id": 1,
      "full_name": "John Doe",
      "email": "user@example.com",
      "role": "user"
    }
  }
  ```
### Delete User
- Endpoint: DELETE /me
- Request Headers:
  ```json
  {
    "Authorization": "Bearer <your_access_token>"
  }
  ```
  Response:
  - Status Code: 204 No Content

# Transaction Management Service
- Base path: /api/v1/

## Models
- Transaction
  - Attributes:
    - id: int
    - created_at: datetime
    - updated_at: datetime
    - transaction_time: datetime
    - account_id: int
    - category_id: int
    - tag_ids: list[int]
    - status: str ["pending", "completed", "failed"]
    - transaction_type: str ["income", "expense", "transfer"]
    - payment_method: str ["cash", "card", "bank_transfer", "online_payment"]
    - amount: decimal
    - description: str
    - currency: str

- Tag
  - Attributes:
    - id: int
    - created_at: datetime
    - updated_at: datetime
    - name: str
    - description: str
    - active: bool

## Authentication

Headers:
- X-Account-ID: int (All requests must include this header. All requests are scoped to a specific account.)
- X-User-Id: int (All requests must include this header.)
- X-User-Roles: list[str] (All requests must include this header. Used for authorization.)

## Endpoints

### Transactions

* **POST /transactions/**
    * Description: Creates a new financial transaction.
    * Request Body (JSON):
        ```json
        {
            "transaction_time": "datetime (ISO 8601)",
            "account_id": "integer (ID of the account from the Accounts service)",
            "category_id": "integer (ID of the category from the Categories service)",
            "tag_ids": "[integer] (list of tag IDs managed by this service - optional)",
            "status": "string (enum: 'pending', 'completed', 'failed' - optional)",
            "transaction_type": "string (enum: 'income', 'expense', 'transfer')",
            "payment_method": "string (enum: 'cash', 'card', 'bank_transfer', 'online_payment')",
            "amount": "decimal",
            "description": "string (optional)",
            "currency": "string (ISO 4217)"
        }
        ```
    * Response (201 Created):
        ```json
        {
            "status": "success",
            "data": {
                "id": "integer",
                "created_at": "datetime (ISO 8601)",
                "updated_at": "datetime (ISO 8601)",
                "transaction_time": "datetime (ISO 8601)",
                "account_id": "integer",
                "category_id": "integer",
                "tag_ids": "[integer]",
                "status": "string (optional)",
                "transaction_type": "string",
                "payment_method": "string",
                "amount": "decimal",
                "description": "string (optional)",
                "currency": "string"
            }
        }
        ```
    * Response (400 Bad Request):
        ```json
        {
          "detail": {
            "status": 400,
            "message": "Validation errors."
          }
        }
        ```

* **GET /transactions/**
    * Description: Lists all financial transactions for the authenticated user (with optional query parameters for filtering, ordering, and pagination).
    * Query Parameters (optional): `account_id`, `category_id`, `date_from`, `date_to`, `transaction_type`, `status`, `payment_method`, `tag_id`, `ordering`, `page`, `page_size`.
    * Response (200 OK):
        ```json
        {
            "status": "success",
            "data": {
                "page": "integer",
                "page_size": "integer",
                "total_items": "integer",
                "total_pages": "integer",
                "filters": "object (e.g., {\"category_id\": 1, \"date_gte\": \"2025-04-01\"})",
                "items": [
                    // List of transaction objects (same format as in POST response data)
                ]
            }
        }
        ```

* **GET /transactions/{id}/**
    * Description: Retrieves a specific financial transaction.
    * Response (200 OK):
        ```json
        {
            "status": "success",
            "data": {
                // Transaction object (same format as in POST response data)
            }
        }
        ```
    * Response (404 Not Found):
        ```json
        {
          "detail": {
            "status": 404,
            "message": "Transaction not found."
          }
        }
        ```

* **PUT /transactions/{id}/**
    * Description: Updates a specific financial transaction.
    * Request Body: Similar to POST request body (some fields might be optional).
    * Response (200 OK):
        ```json
        {
            "status": "success",
            "data": {
                // Updated transaction object
            }
        }
        ```
    * Response (400 Bad Request): (Same error format as POST)
    * Response (404 Not Found): (Same error format as GET)

* **DELETE /transactions/{id}/**
    * Description: Deletes a specific financial transaction.
    * Response (204 No Content): Empty response with 204 status code.
    * Response (404 Not Found): (Same error format as GET)

### Tags

* **POST /tags/**
    * Description: Creates a new transaction tag.
    * Request Body:
        ```json
        {
            "name": "string",
            "description": "string (optional)",
            "active": "boolean (optional, default: true)"
        }
        ```
    * Response (201 Created):
        ```json
        {
            "status": "success",
            "data": {
                "id": "integer",
                "created_at": "datetime (ISO 8601)",
                "updated_at": "datetime (ISO 8601)",
                "name": "string",
                "description": "string (optional)",
                "active": "boolean"
            }
        }
        ```
    * Response (400 Bad Request): (Same error format as POST /transactions/)

* **GET /tags/**
    * Description: Lists all transaction tags.
    * Query Parameters (optional): `active`, `ordering`, `page`, `page_size`.
    * Response (200 OK): (Same pagination format as GET /transactions/)

* **GET /tags/{id}/**
    * Description: Retrieves a specific tag.
    * Response (200 OK):
        ```json
        {
            "status": "success",
            "data": {
                // Tag object
            }
        }
        ```
    * Response (404 Not Found): (Same error format as GET /transactions/)

* **PUT /tags/{id}/**
    * Description: Updates a specific tag.
    * Request Body: Similar to POST request body (all fields optional).
    * Response (200 OK):
        ```json
        {
            "status": "success",
            "data": {
                // Updated tag object
            }
        }
        ```
    * Response (400 Bad Request): (Same error format as POST /transactions/)
    * Response (404 Not Found): (Same error format as GET /transactions/)

* **DELETE /tags/{id}/**
    * Description: Deletes a specific tag.
    * Response (204 No Content): Empty response with 204 status code.
    * Response (404 Not Found): (Same error format as GET /transactions/)

# API Gateway
- GraphQL
- HTTP Library: httpx[http2]
- Server: Hypercorn (HTTP/3)
