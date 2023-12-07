package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
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
	logger := log.New(os.Stdout, "[product-api] ", log.LstdFlags)
	storeLogger := log.New(os.Stdout, "[availability-store] ", log.LstdFlags)

	// NoSQL: Initialize Product Repository store
	store, err := domain.New(storeLogger)
	if err != nil {
		logger.Fatal(err)
	}
	store.CreateTables()

	defer store.CloseSession()

	// NoSQL: Checking if the connection was established

	//Initialize the handler and inject said logger
	reservationHandler := handlers.NewReservationsHandler(logger, store)

	//Initialize the router and add a middleware for all the requests
	router := mux.NewRouter()
	router.Use(reservationHandler.MiddlewareContentTypeSet)

	getAvailabilityRouter := router.Methods(http.MethodGet).Subrouter()
	getAvailabilityRouter.HandleFunc("/accomm/{id}/availability", reservationHandler.GetAvailabilityPeriodsByAccommodation)

	getReservationsRouter := router.Methods(http.MethodGet).Subrouter()
	getReservationsRouter.HandleFunc("/availability/{id}/reservations", reservationHandler.GetReservationsByAvailabilityPeriod)

	postAvailabilityRouter := router.Methods(http.MethodPost).Subrouter()
	postAvailabilityRouter.HandleFunc("/accomm/availability", reservationHandler.InsertAvailabilityPeriodByAccommodation)
	postAvailabilityRouter.Use(reservationHandler.MiddlewareAvailabilityPeriodDeserialization)

	postReservationRouter := router.Methods(http.MethodPost).Subrouter()
	postReservationRouter.HandleFunc("/availability/reservations", reservationHandler.InsertReservationByAvailabilityPeriod)
	postReservationRouter.Use(reservationHandler.MiddlewareReservationDeserialization)

	cors := gorillaHandlers.CORS(gorillaHandlers.AllowedOrigins([]string{"*"}))

	//Initialize the server
	server := http.Server{
		Addr:         ":" + port,
		Handler:      cors(router),
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
	if server.Shutdown(timeoutContext) != nil {
		logger.Fatal("Cannot gracefully shutdown...")
	}
	logger.Println("Server stopped")
}
