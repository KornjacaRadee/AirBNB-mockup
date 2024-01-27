package handlers

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"rating_service/domain"
)

type KeyProduct struct{}

type RatingsHandler struct {
	logger *log.Logger
	repo   *domain.RatingsRepo
}

func NewRatingsHandler(l *log.Logger, r *domain.RatingsRepo) *RatingsHandler {
	return &RatingsHandler{l, r}
}

func (r *RatingsHandler) GetHostReservationsByHost(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	hostId := vars["id"]

	ratings, err := r.repo.GetHostRatingsByHost(hostId)
	if err != nil {
		r.logger.Print("Database exception: ", err)
	}

	if ratings == nil {
		return
	}

	err = ratings.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		r.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (r *RatingsHandler) GetHostReservationsByGuest(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	guestId := vars["id"]

	ratings, err := r.repo.GetHostRatingsByGuest(guestId)
	if err != nil {
		r.logger.Print("Database exception: ", err)
	}

	if ratings == nil {
		return
	}

	err = ratings.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		r.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (r *RatingsHandler) GetAccommodationReservationsByAccommodation(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	accommodationId := vars["id"]

	ratings, err := r.repo.GetAccommodationRatingsByAccommodation(accommodationId)
	if err != nil {
		r.logger.Print("Database exception: ", err)
	}

	if ratings == nil {
		return
	}

	err = ratings.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		r.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (r *RatingsHandler) GetAccommodationReservationsByGuest(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	guestId := vars["id"]

	ratings, err := r.repo.GetAccommodationRatingsByGuest(guestId)
	if err != nil {
		r.logger.Print("Database exception: ", err)
	}

	if ratings == nil {
		return
	}

	err = ratings.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		r.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (r *RatingsHandler) GetAccommodationReservationsByHost(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	hostId := vars["id"]

	ratings, err := r.repo.GetAccommodationRatingsByHost(hostId)
	if err != nil {
		r.logger.Print("Database exception: ", err)
	}

	if ratings == nil {
		return
	}

	err = ratings.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		r.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (r *RatingsHandler) InsertHostRating(rw http.ResponseWriter, h *http.Request) {
	rating := h.Context().Value(KeyProduct{}).(*domain.HostRating)

	err := r.repo.InsertHostRating(rating)
	if err != nil {
		r.logger.Print("Database exception: ", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusCreated)
}

func (r *RatingsHandler) InsertAccommodationRating(rw http.ResponseWriter, h *http.Request) {
	rating := h.Context().Value(KeyProduct{}).(*domain.AccommodationRating)

	err := r.repo.InsertAccommodationRating(rating)
	if err != nil {
		r.logger.Print("Database exception: ", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusCreated)
}

func (a *RatingsHandler) MiddlewareHostRatingDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		rating := &domain.HostRating{}
		err := rating.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			a.logger.Fatal(err)
			return
		}

		ctx := context.WithValue(h.Context(), KeyProduct{}, rating)
		h = h.WithContext(ctx)

		next.ServeHTTP(rw, h)
	})
}

func (a *RatingsHandler) MiddlewareAccommodationRatingDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		rating := &domain.AccommodationRating{}
		err := rating.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			a.logger.Fatal(err)
			return
		}

		ctx := context.WithValue(h.Context(), KeyProduct{}, rating)
		h = h.WithContext(ctx)

		next.ServeHTTP(rw, h)
	})
}

func (a *RatingsHandler) MiddlewareContentTypeSet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		a.logger.Println("Method [", h.Method, "] - Hit path :", h.URL.Path)

		rw.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(rw, h)
	})
}

//CHECKER

const jwtSecret = "g3HtH5KZNq3KcWglpIc3eOBHcrxChcY/7bTKG8a5cHtjn2GjTqUaMbxR3DBIr+44"

func getRoleFromToken(tokenString string) (string, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Provide the secret key used to sign the token
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return "", fmt.Errorf("Invalid token: %v", err)
	}

	// Check if the token is valid
	if !token.Valid {
		return "", fmt.Errorf("Invalid token")
	}

	// Extract user role from claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("Invalid token claims")
	}

	// Get user role
	role, ok := claims["roles"].(string)
	if !ok {
		return "", fmt.Errorf("User role not found in token claims")
	}

	return role, nil
}

func getUserIdFromToken(tokenString string) (string, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Provide the secret key used to sign the token
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return "", fmt.Errorf("Invalid token: %v", err)
	}

	// Check if the token is valid
	if !token.Valid {
		return "", fmt.Errorf("Invalid token")
	}

	// Extract user_id from claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("Invalid token claims")
	}

	// Get user_id
	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("User ID not found in token claims")
	}

	return userID, nil
}
