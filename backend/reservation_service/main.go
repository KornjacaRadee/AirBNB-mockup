package main

import (
	"context"
	"github.com/sony/gobreaker"
	"net/http"
	"os"
	"os/signal"
	"reservation_service/client"
	"reservation_service/config"
	"reservation_service/domain"
	"reservation_service/handlers"
	"time"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8081"
	}

	// Initialize context
	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//Initialize the logger we are going to use, with prefix and datetime for every log
	//logger := log.New(os.Stdout, "[product-api] ", log.LstdFlags)
	//storeLogger := log.New(os.Stdout, "[availability-store] ", log.LstdFlags)
	logger := config.NewLogger("./logging/log.log")
	// NoSQL: Initialize Product Repository store

	store, err := domain.New(logger)
	if err != nil {
		logger.Println(err)
	}
	store.CreateTables()

	defer store.CloseSession()

	//Initialize clients for other services

	accommodationClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     10,
		},
	}

	accommodationBreaker := gobreaker.NewCircuitBreaker(
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

	notificationClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     10,
		},
	}

	notificationBreaker := gobreaker.NewCircuitBreaker(
		gobreaker.Settings{
			Name:        "notification",
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

	accommodation := client.NewAccommodationClient(accommodationClient, os.Getenv("ACCOMMODATION_SERVICE_URI"), accommodationBreaker)

	notification := client.NewNotificationClient(notificationClient, os.Getenv("NOTIFICATION_SERVICE_URI"), notificationBreaker)

	//Initialize the handler and inject said logger
	reservationHandler := handlers.NewReservationsHandler(logger, store, accommodation, notification)

	//Initialize the router and add a middleware for all the requests
	router := mux.NewRouter()
	router.Use(reservationHandler.MiddlewareContentTypeSet)

	getAvailabilityRouter := router.Methods(http.MethodGet).Subrouter()
	getAvailabilityRouter.HandleFunc("/accomm/{id}/availability", reservationHandler.GetAvailabilityPeriodsByAccommodation)

	getAvailabilityRouter.HandleFunc("/accomm/{id}/check", reservationHandler.CheckAccommodationForReservations)

	getReservationsRouter := router.Methods(http.MethodGet).Subrouter()
	getReservationsRouter.HandleFunc("/availability/{id}/reservations", reservationHandler.GetReservationsByAvailabilityPeriod)

	getReservationsRouter.HandleFunc("/guest/{id}/reservations", reservationHandler.GetReservationsByGuestId)

	postAvailabilityRouter := router.Methods(http.MethodPost).Subrouter()
	postAvailabilityRouter.HandleFunc("/accomm/availability", reservationHandler.InsertAvailabilityPeriodByAccommodation)
	postAvailabilityRouter.Use(reservationHandler.MiddlewareAvailabilityPeriodDeserialization)

	postReservationRouter := router.Methods(http.MethodPost).Subrouter()
	postReservationRouter.HandleFunc("/availability/reservations", reservationHandler.InsertReservationByAvailabilityPeriod)
	postReservationRouter.Use(reservationHandler.MiddlewareReservationDeserialization)

	deleteRouter := router.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/reservation/delete/{id}", reservationHandler.DeleteReservationByAvailabilityPeriod)

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
			logger.Println(err)
		}
	}()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, os.Kill)

	sig := <-sigCh
	logger.Println("Received terminate, graceful shutdown", sig)

	//Try to shut down gracefully
	if server.Shutdown(timeoutContext) != nil {
		logger.Fatalf("Cannot gracefully shutdown...")
	}
	logger.Println("Server stopped")
}
