# Go Server

An Https Golang server with signup, login, and auto_login routes.

## Installation

Run the main function to start server.

```bash
go run main.go
```

## Testing

Server can be tested with Postman.

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

<ins>Auto Login</ins>
```bash
Method: GET
Url: https://localhost:8080/auto_login
Auth Type: Bearer Token
Token: <token>
```

## Details

* Connection is encrypted through Https with self-signed certificate.
* Signup requires username, password, and email fields.
* Password is hashed and salted with Bcrypt before storage in mongoDB.
* Login and Signup return Auth (token) and Profile (uid, username, email).
* Auto Login validates JWT and returns Profile from uid in token claims. 