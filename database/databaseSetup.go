package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBSet() *mongo.Client {

	context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(context, options.Client().ApplyURI(""))

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context, nil)

	if err != nil {
		log.Println("Connection failed to MongoDB")
		return nil

	}

	fmt.Println("Connection successfully to MongoDB")

	return client

}

var Client *mongo.Client = DBSet()

func UserData(client *mongo.Client, collectionName string) *mongo.Collection {

	var collection *mongo.Collection = client.Database("EcommerceDatabase").Collection(collectionName)
	return collection

}

func ProductData(client *mongo.Client, collectionName string) *mongo.Collection {

	var collection *mongo.Collection = client.Database("EcommerceDatabase").Collection(collectionName)
	return collection

}
