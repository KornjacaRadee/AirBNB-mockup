package handlers

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"rating_service/client"
	"rating_service/domain"
	"time"
)

type KeyProduct struct{}

type RatingsHandler struct {
	logger             *log.Logger
	repo               *domain.RatingsRepo
	reservationClient  client.ReservationClient
	notificationClient client.NotificationClient
}

func NewRatingsHandler(l *log.Logger, r *domain.RatingsRepo, rc client.ReservationClient, nc client.NotificationClient) *RatingsHandler {
	return &RatingsHandler{l, r, rc, nc}
}

func (r *RatingsHandler) GetHostRatingsByHost(rw http.ResponseWriter, h *http.Request) {
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

func (r *RatingsHandler) GetHostRatingsByGuest(rw http.ResponseWriter, h *http.Request) {
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

func (r *RatingsHandler) GetAccommodationRatingsByAccommodation(rw http.ResponseWriter, h *http.Request) {
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

func (r *RatingsHandler) GetAccommodationRatingsByGuest(rw http.ResponseWriter, h *http.Request) {
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

func (r *RatingsHandler) GetAccommodationRatingsByHost(rw http.ResponseWriter, h *http.Request) {
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
	reservations, err := r.reservationClient.GetReservationsByGuestId(h.Context(), rating.GuestId)
	if err != nil {
		r.logger.Print("Cant get reservations: ", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	guestsStayedWithHostInPast := false

	for _, reservation := range reservations {
		if reservation.HostId.Hex() == rating.HostId.Hex() && time.Now().After(reservation.EndDate) {
			guestsStayedWithHostInPast = true
		}
	}

	if guestsStayedWithHostInPast {
		err = r.repo.InsertHostRating(rating)
		if err != nil {
			r.logger.Print("Database exception: ", err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		notification := client.NotificationData{
			Host: client.User{Id: rating.HostId},
			Text: "Your have been rated (by " + rating.GuestId.Hex() + ")",
			Time: time.Now(),
		}
		// Call the notification service and handle fallback logic
		_, err = r.notificationClient.SendReservationNotification(h.Context(), notification)
		if err != nil {
			log.Printf("Error creating notification: %v", err)
			http.Error(rw, "Notification service not available, but rating created", http.StatusCreated)
			return
		}
		rw.WriteHeader(http.StatusCreated)
	} else {
		http.Error(rw, "Guest didn't stay with the host in the past so he can't rate him", http.StatusBadRequest)
	}
}

func (r *RatingsHandler) InsertAccommodationRating(rw http.ResponseWriter, h *http.Request) {
	rating := h.Context().Value(KeyProduct{}).(*domain.AccommodationRating)

	reservations, err := r.reservationClient.GetReservationsByGuestId(h.Context(), rating.GuestId)
	if err != nil {
		r.logger.Print("Cant get reservations: ", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	guestsStayedInAccommodationInPast := false

	for _, reservation := range reservations {
		if reservation.AccommodationId.Hex() == rating.AccommodationId.Hex() && time.Now().After(reservation.EndDate) {
			guestsStayedInAccommodationInPast = true
		}
	}

	if guestsStayedInAccommodationInPast {
		err = r.repo.InsertAccommodationRating(rating)
		if err != nil {
			r.logger.Print("Database exception: ", err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		notification := client.NotificationData{
			Host: client.User{Id: rating.HostId},
			Text: "Your accommodation has been rated (by " + rating.GuestId.Hex() + ")",
			Time: time.Now(),
		}
		// Call the notification service and handle fallback logic
		_, err = r.notificationClient.SendReservationNotification(h.Context(), notification)
		if err != nil {
			log.Printf("Error creating notification: %v", err)
			http.Error(rw, "Notification service not available, but rating created", http.StatusCreated)
			return
		}
		rw.WriteHeader(http.StatusCreated)
	} else {
		http.Error(rw, "Guest didn't stay in the accommodation in the past so he can't rate it", http.StatusBadRequest)
	}

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
