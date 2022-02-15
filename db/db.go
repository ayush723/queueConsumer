package db

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var ctx context.Context
var err error
var Client *mongo.Client

func init() {
	ctx = context.Background()
	Client, err = mongo.Connect(ctx,
		options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = Client.Ping(context.TODO(),
		readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
}
