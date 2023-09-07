package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client interface {
	Database(string) Database
}

type Database interface {
	Collection(string) Collection
	Client() Client
}

type Collection interface {
	InsertOne(context.Context, interface{}) (interface{}, error)
}

type SingleResult interface {
	Decode(interface{}) error
}

type MongoClient struct {
	mongoClient *mongo.Client
}

type MongoDatabase struct {
	db *mongo.Database
}

type MongoCollection struct {
	collection *mongo.Collection
}

type MongoSingleResult struct {
	singleResult *mongo.SingleResult
}

func NewClient(ctx context.Context, connection string) (Client, error) {
	time.Local = time.UTC

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connection))
	if err != nil {
		log.Println("Error occured while connecting to mongo")
	}
	log.Println("database connection established")
	return &MongoClient{mongoClient: client}, err
}

func (client *MongoClient) Database(dbName string) Database {
	db := client.mongoClient.Database(dbName)
	return &MongoDatabase{db}
}

func (db *MongoDatabase) Collection(name string) Collection {
	collection := db.db.Collection(name)
	return &MongoCollection{collection}
}

func (db *MongoDatabase) Client() Client {
	client := db.db.Client()
	return &MongoClient{mongoClient: client}
}

func (collection *MongoCollection) InsertOne(ctx context.Context, document interface{}) (interface{}, error) {
	result, err := collection.collection.InsertOne(ctx, document)
	return result.InsertedID, err
}

func (s *MongoSingleResult) Decode(result interface{}) error {
	return s.singleResult.Decode(result)
}
