package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"log"
	"os"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

func establishConnection(ctx context.Context, endpoint string) driver.Connection {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{endpoint},
	})
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	return conn
}

func establishTLSConnection(ctx context.Context, endpoint string, encodedCA string) driver.Connection {
	caCertificate, err := base64.StdEncoding.DecodeString(encodedCA)
	if err != nil {
		panic(err)
	}

	// Prepare TLS configuration
	tlsConfig := &tls.Config{}
	certpool := x509.NewCertPool()
	if success := certpool.AppendCertsFromPEM(caCertificate); !success {
		panic("Invalid certificate")
	}
	tlsConfig.RootCAs = certpool

	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{endpoint},
		TLSConfig: tlsConfig,
	})
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	return conn
}

func AdbConnect() *driver.Client {
	// Endpoint of the ArangoDB instance / cluster proxy
	endpoint := os.Getenv("DB_ENDPOINT")
	// Root password
	rootPassword := os.Getenv("ROOT_PASSWORD")
	// Base64 encoded CA certificate, provided for Arango Cloud
	encodedCA := os.Getenv("CA_CERT")

	var conn driver.Connection

	if encodedCA == "" {
		conn = establishConnection(context.TODO(), endpoint)
	} else {
		conn = establishTLSConnection(context.TODO(), endpoint, rootPassword)
	}

	// Create client
	opts := driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication("root", rootPassword),
	}
	c, err := driver.NewClient(opts)
	if err != nil {
		panic(err)
	}

	return &c
}

func EnsureDocuments(ctx context.Context, collectionName string, collection any) {
	newC := driver.WithOverwriteMode(ctx, driver.OverwriteModeReplace)

	db, err := adbClient.Database(ctx, "circulate")
	if err != nil {
		log.Fatalf("Could not open the database: %v", err)
	}

	col, err := db.Collection(ctx, collectionName)
	if err != nil {
		log.Fatalf("Could not retrieve collection: %v", err)
	}

	_, sliceErr, err := col.CreateDocuments(newC, collection)
	if sliceErr != nil {
		log.Printf("Error writing slice %v", sliceErr)
	}
	if err != nil {
		log.Fatalf("Error writing data: %v", err)
	}
}

func EnsureEdgeDocuments(ctx context.Context, collectionName string, collection []Edge) {
	db, err := adbClient.Database(ctx, "circulate")
	if err != nil {
		log.Fatalf("Could not open the database: %v", err)
	}

	g, err := db.Graph(context.TODO(), "graph")
	if err != nil {
		log.Fatalf("Could not open the graph: %v", err)
	}

	ec, _, err := g.EdgeCollection(ctx, collectionName)
	if err != nil {
		log.Fatalf("Could not retrieve edge collection: %v", err)
	}

	_, sliceErr, err := ec.CreateDocuments(ctx, collection)
	if sliceErr != nil {
		log.Printf("Error writing slice %v", sliceErr)
	}
	if err != nil {
		log.Fatalf("Error writing data: %v", err)
	}
}
