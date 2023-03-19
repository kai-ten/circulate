package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

var adbClient = *AdbConnect()
var databaseName = "circulate"
var vertexCollections = [...]string{
	"malware", "reporter", "threatActor", "ioc", "threatSource",
}

func HandleRequest() {
	// Create database or retrieve database
	db := ensureDatabase(context.TODO(), databaseName, nil)

	// Create the graph to link the vertices and edges
	g := ensureGraph(context.TODO(), db, "graph", nil)

	// Create edge collections
	ensureEdgeCollection(context.TODO(), g, "shared", []string{"threatSource"}, []string{"ioc"})
	ensureEdgeCollection(context.TODO(), g, "reported_to", []string{"reporter"}, []string{"threatSource"})
	ensureEdgeCollection(context.TODO(), g, "reported_by", []string{"ioc"}, []string{"reporter"})
	ensureEdgeCollection(context.TODO(), g, "came_from", []string{"ioc"}, []string{"threatSource"})
	ensureEdgeCollection(context.TODO(), g, "found", []string{"reporter"}, []string{"ioc"})
	ensureEdgeCollection(context.TODO(), g, "found_by", []string{"ioc"}, []string{"reporter"})
	ensureEdgeCollection(context.TODO(), g, "identified", []string{"reporter"}, []string{"malware"})
	ensureEdgeCollection(context.TODO(), g, "identified_by", []string{"malware"}, []string{"reporter"})
	ensureEdgeCollection(context.TODO(), g, "associates_with", []string{"malware", "ioc"}, []string{"ioc", "malware"})

	// Create vertex collections
	for _, vertex := range vertexCollections {
		ensureCollection(context.TODO(), db, vertex)
	}

	log.Println("Database, graph, and collections have been created. Check out your ArangoDB :)")
}

func main() {
	lambda.Start(HandleRequest)
}
