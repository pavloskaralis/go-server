# Go Server

An Https Golang server with signup, login, and autologin routes.

## Installation

Run the main function to start server.

```bash
go run main.go
```

## Testing

Server can be tested with Postman.

Signup
```bash
Method: POST
Url: https://localhost:8080/signup
Body: {
    "username": "<username>",
    "password": "<password>",
    "email": "<email>"
}
```
Login
```bash
Method: POST
Url: https://localhost:8080/login
Body: {
    "username": "<username>",
    "password": "<password>"
}
```

Auto Login
```bash
Method: GET
Url: https://localhost:8080/auto_login
Auth Type: Bearer Token
Token: <token>
```

## Explanation

* Server has 3 routes: /signup, /login, and /auto_login.
* Connection is encrypted through Https with self-signed certificate.
* Signup requires username, password, and email fields.
* Password is hashed and salted with Bcrypt before storage in mongoDB.
* Login and Signup return Auth (token) and Profile (uid, username, email).
* Auto Login validates JWT and returns Profile from uid in token claims. 