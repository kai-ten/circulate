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
