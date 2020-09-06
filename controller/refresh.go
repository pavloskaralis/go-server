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

	//retrieve request; lreturn error if body != model 
	mapToken := map[string]string{}

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &mapToken)
	if err != nil {
		resErr.Error = err.Error()
		json.NewEncoder(w).Encode(resErr)
		return
	}
	refreshToken := mapToken["refresh"]
   
	//verify the token; return error if expired
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
	   //Make sure that the token method conform to "SigningMethodHMAC"
	   if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		  return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	   }
	   return []byte(os.Getenv("SIGNATURE")), nil
	})
	if err != nil {
		resErr.Error = err.Error()
		json.NewEncoder(w).Encode(resErr)
		return
	}

	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		resErr.Error = err.Error()
		json.NewEncoder(w).Encode(resErr)
		return
	}

	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
				resErr.Error = err.Error()
				json.NewEncoder(w).Encode(resErr)
				return
		}
		userId, ok := claims["user_id"].(string) //convert the interface to string
		if err != nil {
			resErr.Error = err.Error()
			json.NewEncoder(w).Encode(resErr)
			return
		}
		//Delete the previous Refresh Token
		deleted, delErr := auth.DeleteAuth(refreshUuid)
		if delErr != nil || deleted == 0 { //if any goes wrong
			resErr.Error = err.Error()
			json.NewEncoder(w).Encode(resErr)
			return
		}
		//Create new pairs of refresh and access tokens
		ts, err := auth.CreateToken(userId)
		if  err != nil {
			resErr.Error = err.Error()
			json.NewEncoder(w).Encode(resErr)
			return
		}

		//save the tokens metadata to redis
		err = auth.CreateAuth(userId, ts)
		if err != nil {
			resErr.Error = err.Error()
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
		json.NewEncoder(w).Encode(resErr)
		return	
	}
  }