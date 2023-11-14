package main

import (
	"auth-service/authHandlers"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

func main() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://mimiki003:mimiki003@mongodb:27017").
		SetAuth(options.Credential{
			Username: "mimiki003",
			Password: "mimiki003",
		})
	// Replace 'your_username' and 'your_password' with your MongoDB username and password

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

	// Initialize routes
	r := mux.NewRouter()
	r.HandleFunc("/register", authHandlers.HandleRegister(client)).Methods("POST")
	r.HandleFunc("/login", authHandlers.HandleLogin(client)).Methods("POST")
	r.HandleFunc("/users", authHandlers.HandleGetAllUsers(client)).Methods("GET")
	r.HandleFunc("/users/{id}", authHandlers.HandleDeleteUser(client)).Methods("DELETE")

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
