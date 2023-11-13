package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
)

func main() {

	// Load the combined client certificate and key
	cert, err := ioutil.ReadFile("X509-cert-4027413070962155973.pem")
	if err != nil {
		log.Fatal(err)
	}

	// Create TLS configuration with the combined certificate and key
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // Set to false in production
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{cert},
		}},
	}

	// Set client options with TLS configuration
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017").
		SetTLSConfig(tlsConfig)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

}
