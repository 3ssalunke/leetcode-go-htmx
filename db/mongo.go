package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Client interface {
	Database(string) Database
	Ping(context.Context) error
	Disconnect(context context.Context) error
}

type Database interface {
	Collection(string) Collection
	Client() Client
}

type Collection interface {
	FindOne(context.Context, interface{}) SingleResult
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
		log.Fatal(err)
	}
	return &MongoClient{mongoClient: client}, err
}

func (client *MongoClient) Database(dbName string) Database {
	db := client.mongoClient.Database(dbName)
	return &MongoDatabase{db}
}

func (client *MongoClient) Ping(ctx context.Context) error {
	return client.mongoClient.Ping(ctx, &readpref.ReadPref{})
}

func (client *MongoClient) Disconnect(ctx context.Context) error {
	return client.mongoClient.Disconnect(ctx)
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

func (collection *MongoCollection) FindOne(ctx context.Context, filter interface{}) SingleResult {
	singleResult := collection.collection.FindOne(ctx, filter)
	return &MongoSingleResult{singleResult}
}

func (s *MongoSingleResult) Decode(result interface{}) error {
	return s.singleResult.Decode(result)
}
