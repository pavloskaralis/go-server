# Go Server

An Https Golang server with /signup, /login, and /profile routes.

## Requirements

Server requires local installation of mongoDB and Redis (JWT tracking).

## Installation

Run the main function to start server.

```bash
go run main.go
```

## Testing

Server can be tested using Postman:

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

## Details

* Connection is encrypted through Https with self-signed certificate.
* Signup requires username, password, and email fields.
* Password is hashed and salted with Bcrypt before storage in mongoDB.
* Login and Signup initate creation and storage of access and refresh JWTs.
* JWT tokens are tracked by Redis and get deleted after expiration. 
* /login and /signup return Auth (tokens) and Profile(uid, username, email).
* /profile is wrapped in auth middleware that checks access token expiration.
* /profile validates access token and returns Profile via uid in token claims. 