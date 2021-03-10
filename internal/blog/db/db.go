package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	dbName = "blogdb"
	dbURI  = "mongodb://localhost:27017"
)

var _client *Client

type Post struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content`
	Title    string             `bson:"title`
}
type Client struct {
	*mongo.Client
}

type Collection struct {
	*mongo.Collection
}

func GetClient() *Client {
	if _client == nil {
		_client = newClient()
		return _client
	} else {
		return _client
	}
}

func GetCollection(name string) *Collection {
	return &Collection{
		GetClient().Database(dbName).Collection(name),
	}
}

func (c *Collection) SaveOne(ctx context.Context, document interface{}) (string, error) {
	res, saveErr := c.InsertOne(ctx, document)
	if saveErr != nil {
		return "", status.Errorf(codes.Internal, "unable to save %s: %v", c.Name(), saveErr)
	}
	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", status.Errorf(codes.Internal, "unable to save %s:", c.Name())
	}
	return id.Hex(), nil
}

func GracefulDisconnect() {
	if _client != nil {
		fmt.Printf("disconnecting from database...")
		_client.Disconnect(context.TODO())
		fmt.Printf("done.\n")
	} else {
		fmt.Println("no database cleanup needed...")
	}
	_client = nil
}

func newClient() *Client {
	opts := &options.ClientOptions{}
	opts.ApplyURI(dbURI)
	dbClient, dbErr := mongo.NewClient(opts)
	if dbErr != nil {
		log.Fatal(dbErr)
	}
	dbConnErr := dbClient.Connect(context.TODO())
	if dbConnErr != nil {
		log.Fatal(dbConnErr)
	}
	return &Client{
		dbClient,
	}
}
