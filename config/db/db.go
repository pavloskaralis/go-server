package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client

func InitDB() (error) {
	db, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return err
	}
	db.Connect(context.TODO())
	//check connection
	err = db.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}
	DB = db
	return nil
}

func GetUserCollection() (*mongo.Collection, error) {
	//connect to local mongo; return error if initialization fails
	collection := DB.Database("GoServer").Collection("users")
	return collection, nil 
}