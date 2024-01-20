// main.go

package main

import (
	"context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/promeneili1/AirBNB-mockup/clients"
	"github.com/promeneili1/AirBNB-mockup/domain"
	"github.com/promeneili1/AirBNB-mockup/profileHandler"
	"github.com/sony/gobreaker"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
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

	accommodationsClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     10,
		},
	}
	accommodationsBreaker := gobreaker.NewCircuitBreaker(
		gobreaker.Settings{
			Name:        "accommodations",
			MaxRequests: 1,
			Timeout:     10 * time.Second,
			Interval:    0,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				return counts.ConsecutiveFailures > 2
			},
			OnStateChange: func(name string, from, to gobreaker.State) {
				logger.Printf("CB '%s' changed from '%s' to '%s'\n", name, from, to)
			},
			IsSuccessful: func(err error) bool {
				if err == nil {
					return true
				}
				errResp, ok := err.(domain.ErrResp)
				return ok && errResp.StatusCode >= 400 && errResp.StatusCode < 500
			},
		},
	)

	//Initialize clients for other services
	accommodations := clients.NewAccommodationsClient(accommodationsClient, os.Getenv("ACCOMMODATIONS_SERVICE_URI"), accommodationsBreaker)

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
	r.HandleFunc("/new", profileHandler.CreateProfileHandler(client)).Methods("POST")
	r.HandleFunc("/all", profileHandler.GetAllProfilesHandler(client)).Methods("GET")
	r.HandleFunc("/{id}", profileHandler.GetProfileByIDHandler(client)).Methods("GET")
	r.HandleFunc("/u/{email}", profileHandler.GetProfileByEmailHandler(client)).Methods("GET")
	r.HandleFunc("/update/{id}", profileHandler.UpdateProfileHandler(client)).Methods("PUT")
	r.HandleFunc("/delete/{id}", profileHandler.DeleteProfileHandler(client)).Methods("DELETE")
}
