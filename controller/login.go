package controller

import (
	"context"
	"encoding/json"
	"go-server/config/db"
	"go-server/model"
	"io/ioutil"
	"net/http"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	//setup response
	w.Header().Set("Content-Type", "application/json")
	var resErr model.ResponseError

	//retrieve request; lreturn error if body != model 
	var user model.User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	if err != nil {
		resErr.Error = err.Error()
		json.NewEncoder(w).Encode(resErr)
		return
	}

	//retrieve collection; return error if mongo retrieval fails
	collection, err := db.GetDBCollection()
	if err != nil {
		resErr.Error = err.Error()
		json.NewEncoder(w).Encode(resErr)
		return 
	}

	//query for existing user
	var result model.User
	err = collection.FindOne(context.TODO(), bson.D{{"username", user.Username}}).Decode(&result)
	//return error if user not found
	if err != nil {
		resErr.Error = "User not found."
		json.NewEncoder(w).Encode(resErr)
		return
	}
	//hash provided password; check against db hashed pw; return error if not a match
	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))
	if err != nil {
		resErr.Error = "Invalid password."
		json.NewEncoder(w).Encode(resErr)
		return
	}

	//generate token; return error if jwt fails
	uid := result.UID.Hex(); 
	tokenString, err := CreateToken(uid)
	if err != nil {
		resErr.Error = "Error generating token, try again."
		json.NewEncoder(w).Encode(resErr)
		return
	}

	//return auth and profile
	resSuc := model.ResponseSuccess{
		Auth: model.Auth{
			Token: tokenString,
		},
		Profile: model.Profile{
			UID: uid,
			Username: result.Username,
			Email: result.Email,
		},
	}

	json.NewEncoder(w).Encode(resSuc)
	return
}