package controller

import (
	"context"
	"encoding/json"
	"go-server/config/db"
	"go-server/model"
	"io/ioutil"
	"net/http"
	"go-server/config/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

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

	//validate struct; return error if empty field
	if user.Username == "" || user.Email == "" || user.Password == "" {
		switch {
			case user.Username == "" : resErr.Error = "Missing username field."; 
			case user.Email == "" : resErr.Error = "Missing email field."; 
			case user.Password == "" : resErr.Error = "Missing password field."; 
		}
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

	//query for existing user or email
	var result model.User
	err = collection.FindOne(context.TODO(), bson.D{
		{"$or", bson.A{
			bson.D{{"username", user.Username}},
			bson.D{{"email", user.Email}},
		}},
	}).Decode(&result)
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


			//update model pw; store as new user; return error if mongo fails
			user.Password = string(hash)
			newUser, err := collection.InsertOne(context.TODO(), user)
			if err != nil  {
				resErr.Error = "Error creating user, try again."
				json.NewEncoder(w).Encode(resErr)
				return
			}
			oid,_ := newUser.InsertedID.(primitive.ObjectID); 
			uid := oid.Hex()
			//generate token; return error if jwt fails
			ts, err := auth.CreateToken(uid)
			if err != nil {
				collection.DeleteOne(context.TODO(), bson.M{"_id": oid})
				resErr.Error = "Error generating token, try again."
				json.NewEncoder(w).Encode(resErr)
				return
			}
			//create auth; return error if redis fails
			err = auth.CreateAuth(uid, ts)
			if err != nil {
				collection.DeleteOne(context.TODO(), bson.M{"_id": oid})
				resErr.Error = "Error creating auth, try again."
				json.NewEncoder(w).Encode(resErr)
				return
			}
			
			//return auth and profile
			resSuc := ResponseSuccess{
				Auth: Auth{
					Access: ts.AccessToken,
					Refresh: ts.RefreshToken,
				},
				Profile: Profile{
					UID: uid,
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

	//return error if username or email taken
	switch {
		case result.Username == user.Username: resErr.Error = "Username has been taken.";
		case result.Email == user.Email: resErr.Error = "Email has been taken.";

	}
	json.NewEncoder(w).Encode(resErr)
	return
}

