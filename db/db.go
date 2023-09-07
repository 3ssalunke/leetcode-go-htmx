package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/3ssalunke/leetcode-clone/util"
)

func NewMongoDatabase(config util.Config) Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbHost := config.DBHost
	dbPort := config.DBPort

	mongoDbURI := fmt.Sprintf("mongodb://%s:%s", dbHost, dbPort)

	client, err := NewClient(ctx, mongoDbURI)
	if err != nil {
		log.Fatal(client)
	}

	return client.Database(config.DBName)
}
