package main

import (
	"context"
	"log"

	driver "github.com/arangodb/go-driver"
)

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

func ensureGraph(ctx context.Context, db driver.Database, name string, options *driver.CreateGraphOptions) driver.Graph {
	g, err := db.Graph(ctx, name)
	if driver.IsNotFoundGeneral(err) {
		g, err = db.CreateGraphV2(ctx, name, options)
		if err != nil {
			log.Fatalf("Failed to create graph: %v", name)
		}
	} else if err != nil {
		log.Fatalf("Failed to open graph: %v", name)
	}
	return g
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

func ensureEdgeCollection(ctx context.Context, g driver.Graph, collection string, from, to []string) driver.Collection {
	ec, _, err := g.EdgeCollection(ctx, collection)
	if driver.IsNotFoundGeneral(err) {
		ec, err := g.CreateEdgeCollection(ctx, collection, driver.VertexConstraints{From: from, To: to})
		if err != nil {
			log.Fatalf("Failed to create edge collection: %v", err)
		}
		return ec
	} else if err != nil {
		log.Fatalf("Failed to open edge collection: %v", err)
	}
	return ec
}
