package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"go-server/config/db"
	"go-server/model"
	"io/ioutil"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/mongodb/mongo-go-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	//setup response
	w.Header().Set("Contet-Type", "application/json")
	var res model.ResponseResult

	//retrieve request; return error if body != model 
	var user model.user
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	//retrieve collection; return error if mongo retrieval fails
	collection, err := db.GetDBCollection()
	if err !=nil {
		res.Error = err.Error(); 
		json.NewEcoder(w).Encode(res)
		return
	}

	//query for existing user
	var result model.user
	err = collection.FindOne(context.TODO(), bson.D{{"username", user.Username}}).Decode(&result)
	if err != nil {
		//if no user hash password, and create new user 
		if err.Error() == "mongo: no documents in result" {
			//hash password; return error if encryption fails
			hash, error := bcrypt.generateFromPassword([]byte(user.Password), 5)
			if err != nil {
				res.Error = "Error hashing password, try again."
				json.NewEcoder(w).Encode(res)
			}

			//update model pw; store as new user; return error if mongo fails
			user.Password = string(hash)
			_, err = collection.InsertOne(context.TODO(), user)
			if err != nil  {
				res.Error = "Error creating user, try again."
			}

			//return success message
			res.Result = "Signup successful."
			json.NewEncoder(w).Encode(res)
		}
		//return error if mongo query fails
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	//return error if user exists
	res.Error = "Username already exists."
	json.NewEncoder(w).Encode(res)
	return
}

func LoginHanlder(w http.ResponseWrite, r *http.Request) {
	//setup response
	w.Header().set("Content-Type", "application/json")
	var res mode.ResponseResult

	//retrieve request; lreturn error if body != model 
	var user model.User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	//retrieve collection; return error if mongo retrieval fails
	collection, err := db.GetDBCollection()
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return 
	}

	//query for existing user
	var result model.User
	err = collection.FindOne(context.TODO(), bson.D{{"username", user.Username}}).Decode(&result)
	//return error if user not found
	if err != nil {
		res.Error = "User not found."
		json.NewEncoder(w).Encode(res)
		return
	}
	//hash provided password; check against db hashed pw; return error if not a match
	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))
	if err != nil {
		res.Error = "Invalid password."
		json.NewEncoder(w).Encode(res)
		return
	}

	//generate token; return error if jwt fails
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":  result.Username,
		"email": result.Email,
	})
	tokenString, err := token.SignedString([]byte("aX13bD6u7w2QvGL0"))
	if err != nil {
		res.Error = "Error generating token, try again."
		json.NewEncoder(w).Encode(res)
		return
	}

	//return user profile and token
	result.Token = tokenString
	result.Password =  "" 
	json.NewEncoder(w).Encode(result)
	
}