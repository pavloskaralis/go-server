package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetDBCollection() (*mongo.Collection, error) {
	//connect to local mongo; return error if initialization fails
    client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}
	client.Connect(context.TODO())
	//check connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	client.Connect(context.TODO())
	collection := client.Database("GoServer").Collection("users")
	return collection, nil 
}