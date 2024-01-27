package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rating_service/domain"
	"rating_service/handlers"
	"time"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8086"
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

	//Initialize clients for other services

	//Initialize the handler and inject said logger
	ratingHandler := handlers.NewRatingsHandler(logger, store)

	//Initialize the router and add a middleware for all the requests
	router := mux.NewRouter()
	router.Use(ratingHandler.MiddlewareContentTypeSet)

	getHostRatingsRouter := router.Methods(http.MethodGet).Subrouter()
	getHostRatingsRouter.HandleFunc("/host/{id}/host-ratings", ratingHandler.GetHostReservationsByHost)

	getHostRatingsRouter.HandleFunc("/guest/{id}/host-ratings", ratingHandler.GetHostReservationsByGuest)

	getAccommodationRatingsRouter := router.Methods(http.MethodGet).Subrouter()
	getAccommodationRatingsRouter.HandleFunc("/accommodation/{id}/accommodation-ratings", ratingHandler.GetAccommodationReservationsByAccommodation)

	getAccommodationRatingsRouter.HandleFunc("/host/{id}/accommodation-ratings", ratingHandler.GetAccommodationReservationsByHost)

	getAccommodationRatingsRouter.HandleFunc("/guest/{id}/accommodation-ratings", ratingHandler.GetAccommodationReservationsByGuest)

	postHostRatingRouter := router.Methods(http.MethodPost).Subrouter()
	postHostRatingRouter.HandleFunc("/host-rating", ratingHandler.InsertHostRating)
	postHostRatingRouter.Use(ratingHandler.MiddlewareHostRatingDeserialization)

	postAccommodationRatingRouter := router.Methods(http.MethodPost).Subrouter()
	postAccommodationRatingRouter.HandleFunc("/accommodation-rating", ratingHandler.InsertAccommodationRating)
	postAccommodationRatingRouter.Use(ratingHandler.MiddlewareAccommodationRatingDeserialization)

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
