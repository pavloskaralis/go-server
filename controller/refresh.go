package controller

import (
	"encoding/json"
	"go-server/config/auth"
	"net/http"
	jwt "github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"fmt"
	"os"
)

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	//setup response
	w.Header().Set("Content-Type", "application/json")
	var resErr ResponseError

	//retrieve request; return error if body != model 
	mapToken := map[string]string{}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &mapToken)
	if err != nil {
		resErr.Error = "Invalid request body."
		json.NewEncoder(w).Encode(resErr)
		return
	}
	//check if token provided in request; return error if missing
	refreshToken := mapToken["refresh_token"]
	if len(refreshToken) == 0 {
		resErr.Error = "No refresh token provided."
		json.NewEncoder(w).Encode(resErr)
		return
	}

	//verify signing method; return error if invalid
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
	   //check token method conforms to "SigningMethodHMAC"
	   if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		  return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	   }
	   return []byte(os.Getenv("SIGNATURE")), nil
	})
	if err != nil {
		resErr.Error = "Refresh token is invalid."
		json.NewEncoder(w).Encode(resErr)
		return
	}

	//verify token not expired 
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		resErr.Error = "Refresh token has expired."
		json.NewEncoder(w).Encode(resErr)
		return
	}

	//retrieve claims
	claims, ok := token.Claims.(jwt.MapClaims) 
	if ok && token.Valid {
		//return error if uuid is missing
		refreshUuid, ok := claims["refresh_uuid"].(string)
		if !ok {
			resErr.Error = "Refresh token is missing unique identifer."
			json.NewEncoder(w).Encode(resErr)
			return
		}
		//return error if user id is missing
		userId, ok := claims["uid"].(string) 
		if !ok {
			resErr.Error = "Refresh token is missing user id."
			json.NewEncoder(w).Encode(resErr)
			return
		}
		//delete previous token; return error if redis fails
		deleted, delErr := auth.DeleteAuth(refreshUuid)
		if delErr != nil || deleted == 0 { //if any goes wrong
			resErr.Error = "Refresh token no longer exists."
			json.NewEncoder(w).Encode(resErr)
			return
		}
		//Create new pairs of refresh and access tokens; return error if jwt fails
		ts, err := auth.CreateToken(userId)
		if  err != nil {
			resErr.Error = "Failed to create new auth tokens."
			json.NewEncoder(w).Encode(resErr)
			return
		}
		//save the tokens metadata to redis; return error if redis fails
		err = auth.CreateAuth(userId, ts)
		if err != nil {
			resErr.Error = "Failed to store new auth tokens."
			json.NewEncoder(w).Encode(resErr)
			return
		}
		//return auth and profile
		resSuc := Auth{
			Access: ts.AccessToken,
			Refresh: ts.RefreshToken,
		}
		json.NewEncoder(w).Encode(resSuc)
		return
	} else {
		resErr.Error = "Refresh token is missing unique identifer."
		json.NewEncoder(w).Encode(resErr)
		return	
	}
}