package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

var dbClient *mongo.Client

func GetDBClient() *mongo.Client {
	if dbClient != nil {
		return dbClient
	}

	dbUri := os.Getenv("MONGO_DB_URL")
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(dbUri).SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dbClient, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	return dbClient
}

func GetCollection(collection string) *mongo.Collection {
	return GetDBClient().Database("uviedb").Collection(collection)
}
