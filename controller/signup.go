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
		resErr.Error = "Invalid request body."
		json.NewEncoder(w).Encode(resErr)
		return
	}

	//validate User struct; return error missing field
	if user.Username == "" || user.Email == "" || user.Password == "" {
		switch {
			case user.Username == "" : resErr.Error = "Missing username field."; 
			case user.Password == "" : resErr.Error = "Missing password field."; 
			case user.Email == "" : resErr.Error = "Missing email field."; 
		}
		json.NewEncoder(w).Encode(resErr)
		return
	}
	
	//retrieve user collection
	collection := db.GetUserCollection()
	
	//query for existing user or email
	var result model.User
	err = collection.FindOne(context.TODO(), bson.D{
		{"$or", bson.A{
			bson.D{{"username", user.Username}},
			bson.D{{"email", user.Email}},
		}},
	}).Decode(&result)
	if err != nil {
		//if no user found, hash password and create new user 
		if err.Error() == "mongo: no documents in result" {
			
			//hash password; return error if encryption fails
			hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)
			if err != nil {
				resErr.Error = "Error hashing password, try again."
				json.NewEncoder(w).Encode(resErr)
				return
			}
			//set hash as new user password
			user.Password = string(hash)

			//create user doc id
			uid := primitive.NewObjectID()
			user.UID = uid
			stringUID := uid.Hex()

			//generate tokens; return error if jwt fails
			ts, err := auth.CreateToken(stringUID)
			if err != nil {
				resErr.Error = "Error generating token, try again."
				json.NewEncoder(w).Encode(resErr)
				return
			}

			//create auth; return error if redis fails
			err = auth.CreateAuth(stringUID, ts)
			if err != nil {
				resErr.Error = "Error creating auth, try again."
				json.NewEncoder(w).Encode(resErr)
				return
			}

			//store new user; return error if mongo fails
			_, err = collection.InsertOne(context.TODO(), user)
			if err != nil  {
				resErr.Error = "Error creating user, try again."
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
					UID: stringUID,
					Username: user.Username,
					Email: user.Email,
				},
			}
			json.NewEncoder(w).Encode(resSuc)
			return
		}
		//return error if mongo fails for other reason
		resErr.Error = "Error querying database, try again."
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

