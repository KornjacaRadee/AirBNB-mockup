package main

import (
	"context"
	"github.com/sony/gobreaker"
	"log"
	"net/http"
	"notification_service/client"
	"notification_service/domain"
	handlers "notification_service/handlers"
	"os"
	"os/signal"
	"time"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8085"
	}

	// Initialize context
	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//Initialize the logger we are going to use, with prefix and datetime for every log
	logger := log.New(os.Stdout, "[notification-api] ", log.LstdFlags)
	storeLogger := log.New(os.Stdout, "[notification-store] ", log.LstdFlags)

	// NoSQL: Initialize Accommodation Repository store
	store, err := domain.NewNotificationsRepo(storeLogger)
	if err != nil {
		logger.Fatal(err)
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

	profileServiceURI := os.Getenv("PROFILE_SERVICE_URI")
	profile := client.NewProfileClient(profileClient, profileServiceURI, profileBreaker)

	//Initialize the handler and inject logger
	notificationsHandler := handlers.NewNotificationsHandler(logger, store, profile)

	//Initialize the router and add a middleware for all the requests
	router := mux.NewRouter()
	router.Use(notificationsHandler.MiddlewareContentTypeSet)

	getRouter := router.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/all", notificationsHandler.GetAllNotifications)

	postRouter := router.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/new", notificationsHandler.PostNotification)
	postRouter.Use(notificationsHandler.MiddlewareNotificationDeserialization)

	router.HandleFunc("/user-notifications", notificationsHandler.GetUserNotifications).Methods("GET")

	// ...

	// Start server

	//router.Use(handlers2.AuthMiddleware)

	cors := gorillaHandlers.CORS(gorillaHandlers.AllowedOrigins([]string{"*"}))

	//Initialize the server
	server := http.Server{
		Addr:         ":" + port,
		Handler:      cors(router),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
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
	if server.Shutdown(timeoutContext) != nil {
		logger.Fatal("Cannot gracefully shutdown...")
	}
	logger.Println("Server stopped")
}
