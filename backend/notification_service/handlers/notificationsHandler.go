package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"notification_service/client"
	"notification_service/domain"
	"strings"
)

type KeyProduct struct{}

type NotificationsHandler struct {
	logger        *log.Logger
	repo          *domain.NotificationRepo
	profileClient client.ProfileClient
}

// NewNotificationsHandler Injecting the logger makes this code much more testable.
func NewNotificationsHandler(l *log.Logger, r *domain.NotificationRepo, pc client.ProfileClient) *NotificationsHandler {
	return &NotificationsHandler{l, r, pc}
}

func (a *NotificationsHandler) GetAllNotifications(rw http.ResponseWriter, h *http.Request) {

	notifications, err := a.repo.GetAll()
	if err != nil {
		a.logger.Print("Database exception: ", err)
	}

	if notifications == nil {
		return
	}

	err = notifications.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		a.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (a *NotificationsHandler) GetNotification(rw http.ResponseWriter, h *http.Request) {

	vars := mux.Vars(h)
	id := vars["id"]

	notification, err := a.repo.GetByID(id)
	if err != nil {
		a.logger.Print("Database exception: ", err)
	}

	if notification.Id.Hex() != id {
		http.Error(rw, "Notification not found", 404)
		return
	}

	err = notification.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		a.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (a *NotificationsHandler) PostNotification(rw http.ResponseWriter, h *http.Request) {
	notification := h.Context().Value(KeyProduct{}).(*domain.Notification)

	// Insert the notification
	erra := a.repo.Insert(notification)
	if erra != nil {
		http.Error(rw, "Unable to post notification", http.StatusBadRequest)
		a.logger.Fatal(erra)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func (a *NotificationsHandler) DeleteNotification(rw http.ResponseWriter, h *http.Request) {
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

	// Check if the user has the required role or is the owner of the notification
	if role != "host" {
		http.Error(rw, "Unauthorized: Insufficient privileges", http.StatusUnauthorized)
		return
	}

	// Provjeri da li je korisnik vlasnik notifikacije
	notification, err := a.repo.GetByID(id) // Use the new GetByID function
	if err != nil {
		http.Error(rw, "Error getting notification", http.StatusInternalServerError)
		a.logger.Fatal(err)
		return
	}

	idUser, _ := primitive.ObjectIDFromHex(userID)

	if notification.Host.Id != idUser {
		http.Error(rw, "Unauthorized: User is not the owner of the notification", http.StatusUnauthorized)
		return
	}

	a.repo.Delete(id)
	rw.WriteHeader(http.StatusOK)
}

func (a *NotificationsHandler) GetUserNotifications(rw http.ResponseWriter, h *http.Request) {
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
	// Get notifications for the user
	notifications, err := a.repo.GetNotificationsByUserID(userID)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Error getting notifications: %v", err), http.StatusInternalServerError)
		return
	}

	user, err := a.profileClient.GetAllInformationsByUserID(h.Context(), userID)
	if err != nil {
		a.logger.Println("Failed to get HostID from username:", err)
		http.Error(rw, "Failed to get HostID from username", http.StatusBadRequest)
		return
	}

	a.logger.Println(user.Email)

	// Return the notifications as JSON
	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(notifications); err != nil {
		http.Error(rw, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}

}

func (a *NotificationsHandler) MiddlewareNotificationDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		notification := &domain.Notification{}
		err := notification.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			a.logger.Fatal(err)
			return
		}

		ctx := context.WithValue(h.Context(), KeyProduct{}, notification)
		h = h.WithContext(ctx)

		next.ServeHTTP(rw, h)
	})
}

func (a *NotificationsHandler) MiddlewareContentTypeSet(next http.Handler) http.Handler {
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
