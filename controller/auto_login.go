package controller

import (
	"context"
	"encoding/json"
	"go-server/config/db"
	"go-server/model"
	"net/http"
	"strings"
	"fmt"
	"os"
	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AutoLoginHandler(w http.ResponseWriter, r *http.Request) {
	//setup response
	w.Header().Set("Content-Type", "application/json")
	var resErr model.ResponseError

	//validate token
	tokenString := r.Header.Get("Authorization")
	tokenString = strings.Split(tokenString, " ")[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method.")
		}
		return []byte(os.Getenv("SIGNATURE")), nil
	})
	if err != nil {
		resErr.Error = err.Error()
		json.NewEncoder(w).Encode(resErr)
		return 
	}

	//retrieve claims; return error if failure
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//retrieve collection; return error if mongo fails
		collection, err := db.GetDBCollection()
		if err != nil {
			resErr.Error = err.Error()
			json.NewEncoder(w).Encode(resErr)
			return 
		}
		//search for user by uid; return error if not found
		var result model.User
		uid := claims["uid"].(string)
		objID, _ := primitive.ObjectIDFromHex(uid)
		err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&result)
=		if err != nil {
			resErr.Error = uid
			json.NewEncoder(w).Encode(resErr)
			return
		}
		//return auth and profile
		resSuc := model.ResponseSuccess{
			Auth: model.Auth{
				Token: tokenString,
			},
			Profile: model.Profile{
				UID: result.UID.Hex(),
				Username: result.Username,
				Email: result.Email,
			},
		}
		json.NewEncoder(w).Encode(resSuc)
		return
	} else {
		resErr.Error = err.Error()
		json.NewEncoder(w).Encode(resErr)
		return
	}
}