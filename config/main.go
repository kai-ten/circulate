package main

import (
	"context"
	"log"

	"github.com/arangodb/go-driver"
	"github.com/aws/aws-lambda-go/lambda"
)

var adbClient = *AdbConnect()
var databaseName = "circulate"
var vertexCollections = [...]string{
	"malware", "reporter", "threatActor", "ioc", "threatSource",
}
var edgeCollections = [...]string{
	"shared", "reported_to", "came_from", "found", "found_by", "identified", "associates_with",
}

func ensureDatabase(ctx context.Context, name string, options *driver.CreateDatabaseOptions) driver.Database {
	var db driver.Database

	db, err := adbClient.Database(ctx, name)
	if driver.IsNotFoundGeneral(err) {
		db, err = adbClient.CreateDatabase(ctx, name, options)
		if err != nil {
			if driver.IsConflict(err) {
				log.Fatalf("Failed to create database (conflict): %v", err)
			} else {
				log.Fatalf("Failed to create database: %v", err)
			}
		}
	} else if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	return db
}

func ensureCollection(ctx context.Context, db driver.Database, name string) {
	_, err := db.Collection(context.TODO(), name)
	if driver.IsNotFoundGeneral(err) {
		_, e := db.CreateCollection(context.Background(), name, nil)
		if e != nil {
			log.Fatalf("Failed to create collection: %v", e)
		}
	} else if err != nil {
		log.Fatalf("Failed to open collection: %v", err)
	}
}

func HandleRequest() {
	// Create database or retrieve database
	db := ensureDatabase(context.TODO(), databaseName, nil)

	// Create collections
	for _, vertex := range vertexCollections {
		ensureCollection(context.TODO(), db, vertex)
	}
	for _, edge := range edgeCollections {
		ensureCollection(context.TODO(), db, edge)
	}

	log.Println("Database and collections have been created. Please check ArangoDB.")
}

func main() {
	lambda.Start(HandleRequest)
}
