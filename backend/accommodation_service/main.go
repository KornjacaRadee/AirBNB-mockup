package main

import (
	"accommodation_service/cache"
	"accommodation_service/client"
	"accommodation_service/config"
	"accommodation_service/domain"
	handlers "accommodation_service/handlers"
	"accommodation_service/storage"
	"context"
	"github.com/sony/gobreaker"
	"net/http"
	"os"
	"os/signal"
	"time"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	// Initialize context
	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//Initialize the logger we are going to use, with prefix and datetime for every log
	//logger := log.New(os.Stdout, "[product-api] ", log.LstdFlags)
	//storeLogger := log.New(os.Stdout, "[accommodation-store] ", log.LstdFlags)
	//imageStorageLogger := log.New(os.Stdout, "[accommodation-image_storage] ", log.LstdFlags)
	//redisLogger := log.New(os.Stdout, "[accommodation-cache] ", log.LstdFlags)
	logger := config.NewLogger("./logging/log.log")
	// NoSQL: Initialize Accommodation Repository store
	store, err := domain.New(timeoutContext, logger)
	if err != nil {
		logger.Fatalf("Failed initializing Accommodation Repository Store", err)
	}
	defer store.Disconnect(timeoutContext)

	// NoSQL: Checking if the connection was established
	store.Ping()

	// HDFS: Initializing hdfs storage for images
	images, err := storage.New(logger)
	if err != nil {
		logger.Fatalf("Failed to initialize hdfs storage for images", err)
	}

	// Redis: Initializing redis for image caching
	imageCache := cache.New(logger)
	imageCache.Ping()

	defer images.Close()

	_ = images.CreateDirectories()

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

	reservation := client.NewReservationClient(reservationClient, os.Getenv("RESERVATION_SERVICE_URI"), reservationBreaker)

	//Initialize the handler and inject logger
	accommodationsHandler := handlers.NewAccommodationsHandler(logger, store, imageCache, images, reservation)

	//Initialize the router and add a middleware for all the requests
	router := mux.NewRouter()
	router.Use(accommodationsHandler.MiddlewareContentTypeSet)

	getRouter := router.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/all", accommodationsHandler.GetAllAccommodations)
	getRouter.HandleFunc("/accommodation/walk", accommodationsHandler.WalkRoot)
	getRouter.HandleFunc("/{id}", accommodationsHandler.GetAccommodation)

	postRouter := router.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/new", accommodationsHandler.PostAccommodation)
	postRouter.Use(accommodationsHandler.MiddlewareAccommodationDeserialization)

	postAccommodationImagesRouter := router.Methods(http.MethodPost).Subrouter()
	postAccommodationImagesRouter.HandleFunc("/accommodation/images", accommodationsHandler.CreateAccommodationImages)

	getAccommodationImagesRouter := router.Methods(http.MethodGet).Subrouter()
	getAccommodationImagesRouter.HandleFunc("/accommodation/{id}/images", accommodationsHandler.GetAccommodationImages)
	getAccommodationImagesRouter.Use(accommodationsHandler.MiddlewareCacheAllHit)

	patchRouter := router.Methods(http.MethodPatch).Subrouter()
	patchRouter.HandleFunc("/patch/{id}", accommodationsHandler.PatchAccommodation)
	patchRouter.Use(accommodationsHandler.MiddlewareAccommodationDeserialization)

	deleteRouter := router.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/delete/{id}", accommodationsHandler.DeleteAccommodation)

	deleteRouter.HandleFunc("/delete", accommodationsHandler.DeleteUserAccommodations)

	// Add the search endpoint
	router.HandleFunc("/search", accommodationsHandler.SearchAccommodations).Methods("POST")

	router.HandleFunc("/user-accommodations", accommodationsHandler.GetUserAcommodations).Methods("GET")

	// ...

	// Start server

	//router.Use(handlers2.AuthMiddleware)

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
			logger.Fatalf("Cannot distribute all the connections to goroutines", err)
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
