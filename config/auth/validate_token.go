package auth

import (
	"net/http"
	jwt "github.com/dgrijalva/jwt-go"
)


//check token expiration
func ValidateToken(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
	   return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
	   return err
	}
	return nil
}