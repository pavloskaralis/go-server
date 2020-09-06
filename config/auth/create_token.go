package auth

import (
	"time"
	"os"
	jwt "github.com/dgrijalva/jwt-go"
)

// type TokenDetails struct {
// 	AccessToken  string
// 	RefreshToken string
// 	AccessUuid   string
// 	RefreshUuid  string
// 	AtExpires    int64
// 	RtExpires    int64
// }

func CreateToken(userid string) (string, error) {
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["uid"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("SIGNATURE")))
	if err != nil {
	   return "", err
	}
	return token, nil
  }