package controller

import (
	"context"
	"encoding/json"
	"go-server/config/db"
	"go-server/model"
	"io/ioutil"
	"net/http"
	
	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type ResponseError struct {
	Error string `json:"error"`
}

type Auth struct {
	Token string `json:"token"`
} 

type Profile struct {
	UID string `json:"uid"`
	Username string `json:"username"`
	Email string `json:"email"`
} 

type ResponseSuccess struct {
	Auth Auth
	Profile Profile 
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	//setup response
	w.Header().Set("Content-Type", "application/json")
	var resErr ResponseError

	//retrieve request; return error if body != model 
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
	if err !=nil {
		resErr.Error = err.Error(); 
		json.NewEncoder(w).Encode(resErr)
		return
	}

	//query for existing user
	var result model.User
	err = collection.FindOne(context.TODO(), bson.D{{"username", user.Username}}).Decode(&result)
	if err != nil {
		//if no user hash password, and create new user 
		if err.Error() == "mongo: no documents in result" {
			//hash password; return error if encryption fails
			hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)
			if err != nil {
				resErr.Error = "Error hashing password, try again."
				json.NewEncoder(w).Encode(resErr)
				return
			}

			//generate token; return error if jwt fails
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"username":  result.Username,
				"email": result.Email,
			})
			tokenString, err := token.SignedString([]byte("aX13bD6u7w2QvGL0"))
			if err != nil {
				resErr.Error = "Error generating token, try again."
				json.NewEncoder(w).Encode(resErr)
				return
			}

			//update model pw; store as new user; return error if mongo fails
			user.Password = string(hash)
			newUser, err := collection.InsertOne(context.TODO(), user)
			if err != nil  {
				resErr.Error = "Error creating user, try again."
				json.NewEncoder(w).Encode(resErr)
				return
			}
			oid,_ := newUser.InsertedID.(primitive.ObjectID); 
				
			
			//return auth and profile
			resSuc := ResponseSuccess{
				Auth: Auth{
					Token: tokenString,
				},
				Profile: Profile{
					UID: oid.Hex(),
					Username: user.Username,
					Email: user.Email,
				},
			}
			json.NewEncoder(w).Encode(resSuc)
			return
		}
		//return error if mongo query fails
		resErr.Error = err.Error()
		json.NewEncoder(w).Encode(resErr)
		return
	}

	//return error if user exists
	resErr.Error = "Username already exists."
	json.NewEncoder(w).Encode(resErr)
	return
}

// func LoginHandler(w http.ResponseWriter, r *http.Request) {
// 	//setup response
// 	w.Header().Set("Content-Type", "application/json")
// 	var res model.ResponseResult

// 	//retrieve request; lreturn error if body != model 
// 	var user model.User
// 	body, _ := ioutil.ReadAll(r.Body)
// 	err := json.Unmarshal(body, &user)
// 	if err != nil {
// 		res.Error = err.Error()
// 		json.NewEncoder(w).Encode(res)
// 		return
// 	}

// 	//retrieve collection; return error if mongo retrieval fails
// 	collection, err := db.GetDBCollection()
// 	if err != nil {
// 		res.Error = err.Error()
// 		json.NewEncoder(w).Encode(res)
// 		return 
// 	}

// 	//query for existing user
// 	var result model.User
// 	err = collection.FindOne(context.TODO(), bson.D{{"username", user.Username}}).Decode(&result)
// 	//return error if user not found
// 	if err != nil {
// 		res.Error = "User not found."
// 		json.NewEncoder(w).Encode(res)
// 		return
// 	}
	// //hash provided password; check against db hashed pw; return error if not a match
	// err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))
	// if err != nil {
	// 	res.Error = "Invalid password."
	// 	json.NewEncoder(w).Encode(res)
	// 	return
	// }

	// //generate token; return error if jwt fails
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"username":  result.Username,
	// 	"email": result.Email,
	// })
	// tokenString, err := token.SignedString([]byte("aX13bD6u7w2QvGL0"))
	// if err != nil {
	// 	res.Error = "Error generating token, try again."
	// 	json.NewEncoder(w).Encode(res)
	// 	return
	// }

// 	//return user profile and token
// 	result.Token = tokenString
// 	result.Password =  "" 
// 	json.NewEncoder(w).Encode(result)
	
// }