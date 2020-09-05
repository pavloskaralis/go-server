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
				res.Error = "Error While Hashing Password, Try Again"
				json.NewEcoder(w).Encode(res)
			}

			//update model pw; store as new user; return error if mongo fails
			user.Password = string(hash)
			_, err = collection.InsertOne(context.TODO(), user)
			if err != nil  {
				res.Error = "Error While Creating User, Try Again"
			}

			//return success message
			res.Result = "Signup Successful"
			json.NewEncoder(w).Encode(res)
		}
		//return error if mongo query fails
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	//return error if user exists
	res.Error = "Username already Exists!!"
	json.NewEncoder(w).Encode(res)
	return
}