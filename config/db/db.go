package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client

//called from main
func InitDB() (error) {
	//connect to db
	uri := "mongodb://localhost:27017"
	db, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	db.Connect(context.TODO())
	//check connection
	err = db.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}
	//set global var
	DB = db
	return nil
}

func GetUserCollection() (*mongo.Collection) {
	collection := DB.Database("GoServer").Collection("users")
	return collection
}