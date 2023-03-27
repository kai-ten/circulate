package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

type Edge struct {
	Name string
	From []string
	To   []string
}

var adbClient = *AdbConnect()
var databaseName = "circulate"
var vertexCollections = [...]string{
	"malware", "reporter", "threatActor", "ioc", "threatSource",
}
var edgeCollections = [...]Edge{
	{Name: "shared", From: []string{"threatSource"}, To: []string{"ioc"}},
	{Name: "reported_to", From: []string{"reporter"}, To: []string{"threatSource"}},
	{Name: "reported_by", From: []string{"ioc"}, To: []string{"reporter"}},
	{Name: "came_from", From: []string{"ioc"}, To: []string{"threatSource"}},
	{Name: "found", From: []string{"reporter"}, To: []string{"ioc"}},
	{Name: "found_by", From: []string{"ioc"}, To: []string{"reporter"}},
	{Name: "identified", From: []string{"reporter"}, To: []string{"malware"}},
	{Name: "identified_by", From: []string{"malware"}, To: []string{"reporter"}},
	{Name: "associates_with", From: []string{"malware", "ioc"}, To: []string{"ioc", "malware"}},
}

func HandleRequest() {
	// Create database or retrieve database
	db := ensureDatabase(context.TODO(), databaseName, nil)

	// Create the graph to link the vertices and edges
	g := ensureGraph(context.TODO(), db, "threat_graph", nil)

	// Create vertex collections
	for _, vertex := range vertexCollections {
		ensureCollection(context.TODO(), db, vertex)
	}

	// Create edge collections
	for _, edge := range edgeCollections {
		ensureEdgeCollection(context.TODO(), g, edge.Name, edge.From, edge.To)
	}

	log.Println("Database, graph, and collections have been created. Check out your ArangoDB :)")
}

func main() {
	lambda.Start(HandleRequest)
}
