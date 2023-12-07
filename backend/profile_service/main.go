// main.go

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/promeneili1/AirBNB-mockup/profileHandler"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Inicijalizacija MongoDB klijenta
	client, err := initializeMongoDB()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	// Konfiguracija logera
	logger := log.New(os.Stdout, "[profile-api] ", log.LstdFlags)

	// Postavljanje servera
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8084"
	}

	r := mux.NewRouter()
	registerProfileRoutes(r, client) // Registrovanje ruta za profile

	headers := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"http://localhost:4200"})

	handlerWithCORS := handlers.CORS(headers, methods, origins)(r)

	server := http.Server{
		Addr:         ":" + port,
		Handler:      handlerWithCORS,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	logger.Println("Server listening on port", port)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	// Graciozno gašenje
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, os.Kill)

	sig := <-sigCh
	logger.Println("Received terminate, graceful shutdown", sig)

	if err := server.Shutdown(context.TODO()); err != nil {
		logger.Fatal("Cannot gracefully shutdown...")
	}
	logger.Println("Server stopped")
}

func initializeMongoDB() (*mongo.Client, error) {
	dbURI := os.Getenv("MONGO_DB_URI")
	clientOptions := options.Client().ApplyURI(dbURI)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func registerProfileRoutes(r *mux.Router, client *mongo.Client) {
	r.HandleFunc("/profiles/new", profileHandler.CreateProfileHandler(client)).Methods("POST")
	r.HandleFunc("/profiles/all", profileHandler.GetAllProfilesHandler(client)).Methods("GET")
	r.HandleFunc("/profiles/byId/{id}", profileHandler.GetProfileByIDHandler(client)).Methods("GET")
	r.HandleFunc("/profiles/updateById/{id}", profileHandler.UpdateProfileHandler(client)).Methods("PUT")
	r.HandleFunc("/profiles/deleteById/{id}", profileHandler.DeleteProfileHandler(client)).Methods("DELETE")
}
