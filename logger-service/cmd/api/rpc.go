package main

import (
	"context"
	"log"
	"log-service/data"
	"log-service/utilities"

	"go.mongodb.org/mongo-driver/mongo"
)

type RPCServer struct{}

type RPCPayload struct {
	Name string
	Data string
}

var (
	databaseLog   *mongo.Database
	collectionLog *mongo.Collection
)

func initCollLog() *mongo.Collection {
	if databaseLog == nil {
		databaseLog = client.Database("logs")
	}

	if collectionLog == nil {
		collectionLog = databaseLog.Collection("logs")
	}

	return collectionLog
}

func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	coll := initCollLog()
	_, err := coll.InsertOne(context.Background(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: utilities.TimeLocalNow(),
	})
	if err != nil {
		log.Println("error wrting to mongo", err)
		return err
	}

	*resp = "Process payload via RPC" + payload.Name

	return nil
}
