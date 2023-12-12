package main

import (
	"auth-service/authHandlers"
	"auth-service/client"
	"context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// da probam samo komit jedan

func main() {
	// Set client options
	// TREBA MI OVO DA TESTIRAM PA MARE SKIDAJ POSLE

	dburi := os.Getenv("MONGO_DB_URI")
	clientOptions := options.Client().ApplyURI(dburi)

	dbClient, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = dbClient.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	logger := log.New(os.Stdout, "[product-api] ", log.LstdFlags)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8082"
	}

	//Initialize clients for other services
	profileClient := client.NewProfileClient(os.Getenv("PROFILE_SERVICE_URI"))

	r := mux.NewRouter()
	r.HandleFunc("/register", authHandlers.HandleRegister(dbClient, profileClient)).Methods("POST")
	r.HandleFunc("/login", authHandlers.HandleLogin(dbClient)).Methods("POST")
	r.HandleFunc("/users", authHandlers.HandleGetAllUsers(dbClient)).Methods("GET")
	r.HandleFunc("/user", authHandlers.HandleDeleteUser(dbClient)).Methods("DELETE")
	r.HandleFunc("/users/{id}", authHandlers.HandleGetUserByID(dbClient)).Methods("GET")
	// change user passwrod
	r.HandleFunc("/change-password", authHandlers.HandleChangePassword(dbClient)).Methods("POST")

	// initiate password recovery
	r.HandleFunc("/password-recovery", authHandlers.HandlePasswordRecovery(dbClient)).Methods("POST")

	// validate token and lets user access to update password page
	r.HandleFunc("/reset", authHandlers.HandlePasswordReset(dbClient)).Methods("GET")
	// updates users password with new one
	r.HandleFunc("/update", authHandlers.HandlePasswordUpdate(dbClient)).Methods("POST")

	// Enable CORS
	headers := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"http://localhost:4200"}) // Update with your Angular app's origin

	// Apply CORS middleware
	handlerWithCORS := handlers.CORS(headers, methods, origins)(r)

	http.Handle("/", handlerWithCORS)

	//Initialize the server
	server := http.Server{
		Addr:         ":" + port,
		Handler:      handlerWithCORS,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	logger.Println("Server listening on port", port)
	//Distribute all the connections to goroutines
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, os.Kill)

	sig := <-sigCh
	logger.Println("Received terminate, graceful shutdown", sig)

	//Try to shut down gracefully
	if server.Shutdown(context.TODO()) != nil {
		logger.Fatal("Cannot gracefully shutdown...")
	}
	logger.Println("Server stopped")
}
