package main

import (
	"auth-service/data"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"net/http"
)

type UsersHandler struct {
	logger   *log.Logger
	userRepo data.UserDB
}

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

	// Initialize routes
	r := mux.NewRouter()
	r.HandleFunc("/register", handleRegister(client)).Methods("POST")
	r.HandleFunc("/login", handleLogin(client)).Methods("POST")
	r.HandleFunc("/users", handleGetAllUsers(client)).Methods("GET")
	r.HandleFunc("/users/{id}", handleDeleteUser(client)).Methods("DELETE")

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleRegister(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newUser User
		if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Validate user input
		if err := validateUserInput(&newUser); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Register the user
		if err := data.userDB.registerUser(client, &newUser); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func handleLogin(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract email and password from the request
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Login the user
		user, err := loginUser(client, email, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Return user data (you might want to exclude sensitive information like passwords)
		json.NewEncoder(w).Encode(user)
	}
}

func handleGetAllUsers(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := getAllUsers(client)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(users)
	}
}

func handleDeleteUser(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID, err := primitive.ObjectIDFromHex(vars["id"])
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		if err := deleteUser(client, userID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func validateUserInput(user *User) error {
	// Implement validation logic here
	return nil
}
