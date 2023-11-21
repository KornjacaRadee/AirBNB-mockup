package main

import (
	"auth-service/authHandlers"
	"context"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

// da probam samo komit jedan

func main() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://mimiki003:mimiki003@mongodb:27017").
		SetAuth(options.Credential{
			Username: "mimiki003",
			Password: "mimiki003",
		})

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	r := mux.NewRouter()
	r.HandleFunc("/register", authHandlers.HandleRegister(client)).Methods("POST")
	r.HandleFunc("/login", authHandlers.HandleLogin(client)).Methods("POST")
	r.HandleFunc("/users", authHandlers.HandleGetAllUsers(client)).Methods("GET")
	r.HandleFunc("/users/{id}", authHandlers.HandleDeleteUser(client)).Methods("DELETE")

	// Enable CORS
	headers := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"http://localhost:4200"}) // Update with your Angular app's origin

	// Apply CORS middleware
	handlerWithCORS := handlers.CORS(headers, methods, origins)(r)

	http.Handle("/", handlerWithCORS)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
