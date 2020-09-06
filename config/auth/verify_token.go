package auth

import (
	"net/http"
	jwt "github.com/dgrijalva/jwt-go"
	"fmt"
	"os"
)

//verify token signature 
func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	   //Make sure that the token method conform to "SigningMethodHMAC"
	   if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		  return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	   }
	   return []byte(os.Getenv("SIGNATURE")), nil
	})
	if err != nil {
	   return nil, err
	}
	return token, nil
}