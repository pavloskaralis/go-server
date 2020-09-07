package controller

import (
	"context"
	"encoding/json"
	"go-server/config/db"
	"go-server/model"
	"net/http"
	"go-server/config/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	var resErr ResponseError

	tokenAuth, err := auth.ExtractTokenMetadata(r)
	if err != nil {
		resErr.Error = "Error extracting token data, try again."
		json.NewEncoder(w).Encode(resErr)
		return 
	}

	userId, err := auth.FetchAuth(tokenAuth)
	if err != nil {
		resErr.Error = "Error fetching authentication, try again."
		json.NewEncoder(w).Encode(resErr)
		return 
	}

	//retrieve collection
	collection := db.GetUserCollection()
	
	//search for user by uid; return error if not found
	var result model.User
	objID, _ := primitive.ObjectIDFromHex(userId)
	err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&result)
	if err != nil {
		resErr.Error = "User not found."
		json.NewEncoder(w).Encode(resErr)
		return
	}
	
	//return auth and profile
	resSuc := Profile{ 
		UID: result.UID.Hex(),
		Username: result.Username,
		Email: result.Email,
	}
	json.NewEncoder(w).Encode(resSuc)
	return
}