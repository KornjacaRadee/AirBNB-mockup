package handlers

import (
	"accommodation_service/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"strings"
)

type KeyProduct struct{}

type AccommodationsHandler struct {
	logger *log.Logger
	// NoSQL: injecting accommodation repository
	repo *domain.AccommodationRepo
}

// NewAccommodationsHandler Injecting the logger makes this code much more testable.
func NewAccommodationsHandler(l *log.Logger, r *domain.AccommodationRepo) *AccommodationsHandler {
	return &AccommodationsHandler{l, r}
}

func (a *AccommodationsHandler) GetAllAccommodations(rw http.ResponseWriter, h *http.Request) {

	accommodations, err := a.repo.GetAll()
	if err != nil {
		a.logger.Print("Database exception: ", err)
	}

	if accommodations == nil {
		return
	}

	err = accommodations.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		a.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (a *AccommodationsHandler) GetAccommodation(rw http.ResponseWriter, h *http.Request) {

	vars := mux.Vars(h)
	id := vars["id"]

	accommodation, err := a.repo.GetByID(id)
	if err != nil {
		a.logger.Print("Database exception: ", err)
	}

	if accommodation.Id.Hex() != id {
		http.Error(rw, "Accommodation not found", 404)
		return
	}

	err = accommodation.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		a.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (a *AccommodationsHandler) PostAccommodation(rw http.ResponseWriter, h *http.Request) {
	tokenString := h.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(rw, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	// Remove 'Bearer ' prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	role, err := getRoleFromToken(tokenString)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Error extracting user role: %v", err), http.StatusUnauthorized)
		return
	}

	// Check if the user has the required role
	if role != "host" {
		http.Error(rw, "Unauthorized: Insufficient privileges", http.StatusUnauthorized)
		return
	}

	// Extract user ID from the token
	userID, err := getUserIdFromToken(tokenString)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Error extracting user ID: %v", err), http.StatusUnauthorized)
		return
	}

	// Create a new accommodation with the extracted user ID as the owner
	accommodation := h.Context().Value(KeyProduct{}).(*domain.Accommodation)
	accommodation.Owner.Id, _ = primitive.ObjectIDFromHex(userID)

	// Insert the accommodation
	erra := a.repo.Insert(accommodation)
	if erra != nil {
		http.Error(rw, "Unable to post accommodation", http.StatusBadRequest)
		a.logger.Fatal(erra)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func (a *AccommodationsHandler) PatchAccommodation(rw http.ResponseWriter, h *http.Request) {
	tokenString := h.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(rw, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	// Remove 'Bearer ' prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	role, err := getRoleFromToken(tokenString)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Error extracting user role: %v", err), http.StatusUnauthorized)
		return
	}

	// Check if the user has the required role
	if role != "host" {
		http.Error(rw, "Unauthorized: Insufficient privileges", http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(h)
	id := vars["id"]
	accommodation := h.Context().Value(KeyProduct{}).(*domain.Accommodation)

	a.repo.Update(id, accommodation)
	rw.WriteHeader(http.StatusOK)
}

//func (a *AccommodationsHandler) DeleteAccommodation(rw http.ResponseWriter, h *http.Request) {
//
//	vars := mux.Vars(h)
//	id := vars["id"]
//
//	a.repo.Delete(id)
//	rw.WriteHeader(http.StatusNoContent)
//}

func (a *AccommodationsHandler) DeleteAccommodation(rw http.ResponseWriter, h *http.Request) {
	tokenString := h.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(rw, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	// Remove 'Bearer ' prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	role, err := getRoleFromToken(tokenString)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Error extracting user role: %v", err), http.StatusUnauthorized)
		return
	}

	// Extract user ID from the token
	userID, err := getUserIdFromToken(tokenString)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Error extracting user ID: %v", err), http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(h)
	id := vars["id"]

	// Check if the user has the required role or is the owner of the accommodation
	if role != "host" {
		http.Error(rw, "Unauthorized: Insufficient privileges", http.StatusUnauthorized)
		return
	}

	// Provjeri da li je korisnik vlasnik smjestaja
	accommodation, err := a.repo.GetByID(id) // Use the new GetByID function
	if err != nil {
		http.Error(rw, "Error getting accommodation", http.StatusInternalServerError)
		a.logger.Fatal(err)
		return
	}

	if accommodation.Owner.Id.String() != userID {
		http.Error(rw, "Unauthorized: User is not the owner of the accommodation", http.StatusUnauthorized)
		return
	}

	a.repo.Delete(id)
	rw.WriteHeader(http.StatusNoContent)
}

func (a *AccommodationsHandler) GetUserAcommodations(rw http.ResponseWriter, h *http.Request) {
	tokenString := h.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(rw, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	// Remove 'Bearer ' prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	userID, err := getUserIdFromToken(tokenString)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Error extracting user ID: %v", err), http.StatusUnauthorized)
		return
	}
	// Extract user ID from JWT token

	// Get accommodations for the user
	accommodations, err := a.repo.GetAccommodationsByUserID(userID)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Error getting accommodations: %v", err), http.StatusInternalServerError)
		return
	}

	// Return the accommodations as JSON
	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(accommodations); err != nil {
		http.Error(rw, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}

}

func (a *AccommodationsHandler) MiddlewareAccommodationDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		accommodation := &domain.Accommodation{}
		err := accommodation.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			a.logger.Fatal(err)
			return
		}

		ctx := context.WithValue(h.Context(), KeyProduct{}, accommodation)
		h = h.WithContext(ctx)

		next.ServeHTTP(rw, h)
	})
}

func (a *AccommodationsHandler) SearchAccommodations(rw http.ResponseWriter, h *http.Request) {
	var searchRequest domain.SearchRequest
	err := searchRequest.FromJSON(h.Body)
	if err != nil {
		http.Error(rw, "Unable to decode search request", http.StatusBadRequest)
		a.logger.Fatal(err)
		return
	}

	accommodations, err := a.repo.SearchAccommodations(searchRequest)
	if err != nil {
		http.Error(rw, "Error searching accommodations", http.StatusInternalServerError)
		a.logger.Fatal(err)
		return
	}

	err = accommodations.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		a.logger.Fatal("Unable to convert to json:", err)
		return
	}
}

func (a *AccommodationsHandler) MiddlewareContentTypeSet(next http.Handler) http.Handler {
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