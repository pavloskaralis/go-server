package auth

import (
	"net/http"
	jwt "github.com/dgrijalva/jwt-go"
)


type AccessDetails struct {
    AccessUuid string
    UserId   string
}

//extract uid 
func ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
	   return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
	   accessUuid, ok := claims["access_uuid"].(string)
	   userId, ok := claims["uid"].(string)
	   if !ok {
		  return nil, err
	   }
	   if err != nil {
		  return nil, err
	   }
	   return &AccessDetails{
		  AccessUuid: accessUuid,
		  UserId:   userId,
	   }, nil
	}
	return nil, err
  }