package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	grpcPort = "50001"
	monGoURL = "mongodb://mongo:27017"
)

var client *mongo.Client

type Config struct {
}

func main() {
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

}

func connectToMongo() (*mongo.Client, error) {
	opts := options.Client().ApplyURI(mongoURL)
	opts.SetAuth(options.Credential{
		Username: "admin",    // Change this one
		Password: "password", // Change this one
	})

	connection, err := mongo.Connect(context.Backround(), otps)
	if err != nil {
		log.Println("error when connect", err)
		return nil, err
	}

	return c, nil
}
