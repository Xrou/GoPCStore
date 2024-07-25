package model

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type database struct {
	Client   *mongo.Client
	Database *mongo.Database
}

var Database = database{}

func Connect() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017/"))

	if err != nil {
		panic(err)
	}

	db := client.Database("ComputerShop")

	Database.Client = client
	Database.Database = db
}

func Disconnect() {
	Database.Client.Disconnect(context.Background())
}
