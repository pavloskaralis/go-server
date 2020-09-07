# Go Server

An Https Golang server with /signup, /login, /refresh, and /profile routes.

## Requirements

Server requires local copy of mongoDB and Redis (JWT tracking).

Mongo URI: mongodb://localhost:27017
Redis DNS: localhost:6379

## Installation

Run the main function to start server.

```bash
go run main.go
```

## Testing

Server can be tested using Postman.

<ins>Signup</ins>
```bash
Method: POST
Url: https://localhost:8080/signup
Body: {
    "username": "<username>",
    "password": "<password>",
    "email": "<email>"
}
```
<ins>Login</ins>
```bash
Method: POST
Url: https://localhost:8080/login
Body: {
    "username": "<username>",
    "password": "<password>"
}
```

<ins>Profile</ins>
```bash
Method: GET
Url: https://localhost:8080/profile
Auth Type: Bearer Token
Token: <access token>
```

<ins>Refresh</ins>
```bash
Method: POST
Url: https://localhost:8080/refresh
Body: {
    "refresh": "<refresh_token>",
}
## Details

* Connection is encrypted through Https with self-signed certificate.
* Signup requires username, password, and email fields.
* Password is hashed and salted by Bcrypt before storage in mongoDB.
* /login and /signup initiate creation of access and refresh JWTs.
* JWT tokens are tracked by Redis and get deleted after expiration. 
* /login and /signup return Auth (tokens) and Profile(uid, username, email).
* /profile is wrapped in auth middleware that checks access token expiration.
* /profile validates access token and returns Profile via uid in token claims. 
* /refresh returns refreshed Auth if provided a valid refresh token.