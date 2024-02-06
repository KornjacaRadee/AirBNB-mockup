package main

import (
	"auth-service/authHandlers"
	"auth-service/client"
	"auth-service/config"
	"auth-service/domain"
	"context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/sony/gobreaker"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	logger := config.NewLogger("./logging/log.log")

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8082"
	}
	profileClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     10,
		},
	}
	profileBreaker := gobreaker.NewCircuitBreaker(
		gobreaker.Settings{
			Name:        "profile",
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

	reservationClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     10,
		},
	}
	reservationBreaker := gobreaker.NewCircuitBreaker(
		gobreaker.Settings{
			Name:        "reservation",
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

	accommClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     10,
		},
	}
	accommBreaker := gobreaker.NewCircuitBreaker(
		gobreaker.Settings{
			Name:        "accommodation",
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
	profile := client.NewProfileClient(profileClient, os.Getenv("PROFILE_SERVICE_URI"), profileBreaker)
	reservation := client.NewReservationClient(reservationClient, os.Getenv("RESERVATION_SERVICE_URI"), reservationBreaker)
	accommodation := client.NewAccommodationClient(accommClient, os.Getenv("ACCOMMODATION_SERVICE_URI"), accommBreaker)

	r := mux.NewRouter()
	r.HandleFunc("/register", authHandlers.HandleRegister(dbClient, profile)).Methods("POST")
	r.HandleFunc("/login", authHandlers.HandleLogin(dbClient)).Methods("POST")
	r.HandleFunc("/users", authHandlers.HandleGetAllUsers(dbClient)).Methods("GET")
	r.HandleFunc("/user", authHandlers.HandleDeleteUser(dbClient, reservation, accommodation, profile)).Methods("DELETE")
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
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Println("Server listening on port", port)
	//Distribute all the connections to goroutines
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Panicf("Panic on auth-service during listening")
		}
	}()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, os.Kill)

	sig := <-sigCh
	logger.Println("Received terminate, graceful shutdown", sig)

	//Try to shut down gracefully
	if server.Shutdown(context.TODO()) != nil {
		logger.Fatalf("Cannot gracefully shutdown...")
	}
	logger.Println("Server stopped")
}
