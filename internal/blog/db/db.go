package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BlogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content`
	Title    string             `bson:"title`
}

func NewClient() *mongo.Client {
	opts := &options.ClientOptions{
		Hosts: []string{"mongodb://localhost:27017"},
	}
	dbClient, dbErr := mongo.NewClient(opts)
	if dbErr != nil {
		log.Fatal(dbErr)
	}
	dbConnErr := dbClient.Connect(context.TODO())
	if dbConnErr != nil {
		log.Fatal(dbConnErr)
	}
	return dbClient
}
