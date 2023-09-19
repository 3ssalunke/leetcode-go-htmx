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

type Cursor interface {
	All(context.Context, interface{}) error
}

type Collection interface {
	Find(context.Context, interface{}, ...*options.FindOptions) (Cursor, error)
	FindOne(context.Context, interface{}) SingleResult
	InsertOne(context.Context, interface{}) (interface{}, error)
	Aggregate(context.Context, interface{}, ...*options.AggregateOptions) (Cursor, error)
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

type MongoCursor struct {
	cursor *mongo.Cursor
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

func (collection *MongoCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (Cursor, error) {
	cursor, err := collection.collection.Find(ctx, filter, opts...)
	return &MongoCursor{cursor}, err
}

func (collection *MongoCollection) Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (Cursor, error) {
	cursor, err := collection.collection.Aggregate(ctx, pipeline, opts...)
	return &MongoCursor{cursor}, err
}

func (collection *MongoCollection) FindOne(ctx context.Context, filter interface{}) SingleResult {
	singleResult := collection.collection.FindOne(ctx, filter)
	return &MongoSingleResult{singleResult}
}

func (s *MongoSingleResult) Decode(result interface{}) error {
	return s.singleResult.Decode(result)
}

func (cursor *MongoCursor) All(ctx context.Context, result interface{}) error {
	return cursor.cursor.All(ctx, result)
}
