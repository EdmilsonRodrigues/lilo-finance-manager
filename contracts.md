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

# API Gateway
- GraphQL
- HTTP Library: httpx[http2]
- Server: Hypercorn (HTTP/3)


