package authHandlers

import (
	"auth-service/data"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func HandleRegister(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newUser data.User
		if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Validate user input
		if err := validateUserInput(&newUser); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Hash the password
		hashedPassword, err := data.HashPassword(newUser.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		newUser.Password = hashedPassword

		// Register the user
		if err := data.RegisterUser(client, &newUser); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

const jwtSecret = "g3HtH5KZNq3KcWglpIc3eOBHcrxChcY/7bTKG8a5cHtjn2GjTqUaMbxR3DBIr+44"

func generateJWTToken(user *data.User) (string, error) {
	// Create a new token object, specifying signing method and claims
	claims := jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expiration time (e.g., 24 hours)
		"roles":   user.Roles,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func HandleLogin(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode the JSON request body
		var credentials struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}
		log.Printf("Login request: email=%s, password=%s\n", credentials.Email, credentials.Password)

		// Login the user
		user, err := data.LoginUser(client, credentials.Email, credentials.Password)
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		// Generate a JWT token
		token, err := generateJWTToken(user)
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		// Return user data and token
		response := struct {
			User  *data.User `json:"user"`
			Token string     `json:"token"`
		}{
			User:  user,
			Token: token,
		}

		// Set the "Content-Type" header to "application/json"
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(response)
	}
}

// I THINK THIS FUNC SHOULD NOT BE AVAILABLE TO REQUEST

//func HandleGetUserByID(client *mongo.Client) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		// Extract user ID from URL parameters
//		vars := mux.Vars(r)
//		userID, err := primitive.ObjectIDFromHex(vars["id"])
//		if err != nil {
//			http.Error(w, "Invalid user ID", http.StatusBadRequest)
//			return
//		}
//
//		// Get user by ID
//		user, err := data.GetUserByID(client, userID)
//		if err != nil {
//			http.Error(w, "User not found", http.StatusNotFound)
//			return
//		}
//
//		// Return user data
//		w.Header().Set("Content-Type", "application/json")
//		json.NewEncoder(w).Encode(user)
//	}
//}

// authHandlers/authHandlers.go

// ...

func HandleChangePassword(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID from JWT token
		userIDFromToken, err := extractUserIDFromToken(r)
		if err != nil {
			log.Printf("Invalid token: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Convert userIDFromToken to ObjectID
		userID, err := primitive.ObjectIDFromHex(userIDFromToken)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Decode the JSON request body
		rawRequestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading raw request body: %v", err)
			http.Error(w, "Error reading raw request body", http.StatusInternalServerError)
			return
		}
		log.Printf("Raw request body: %s", rawRequestBody)

		// Decode the JSON request body
		var passwordChange struct {
			OldPassword string `json:"old_password"`
			NewPassword string `json:"new_password"`
		}
		if err := json.Unmarshal(rawRequestBody, &passwordChange); err != nil {
			log.Printf("Error decoding JSON: %v", err)
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}
		log.Printf("New password from handler is: %s", passwordChange.NewPassword)

		user, err := validateOldPasswordAndGetUser(client, userID, passwordChange.OldPassword)
		if err != nil {
			log.Printf("Error validating old password: %v", err)
			http.Error(w, "Error validating old password", http.StatusInternalServerError)
			return
		}

		if user == nil {
			http.Error(w, "Invalid old password", http.StatusUnauthorized)
			return
		}

		// Update the user's password in the database
		if err := data.UpdatePassword(client, userID, passwordChange.NewPassword); err != nil {
			log.Printf("Error updating password: %v", err)
			http.Error(w, "Error updating password", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// Helper function to validate the old password and get the user
func validateOldPasswordAndGetUser(client *mongo.Client, userID primitive.ObjectID, oldPassword string) (*data.User, error) {
	user, err := data.GetUserByID(client, userID)
	if err != nil {
		return nil, err
	}

	// Compare the old password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return nil, nil
	}

	return user, nil
}

func HandleGetAllUsers(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := data.GetAllUsers(client)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(users)
	}
}

func HandleDeleteUser(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID from JWT token
		userIDFromToken, err := extractUserIDFromToken(r)
		if err != nil {
			log.Printf("Invalid token: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		log.Printf("Received token: %s", r.Header.Get("Authorization"))

		// Convert user ID from string to primitive.ObjectID
		objectIDFromToken, err := primitive.ObjectIDFromHex(userIDFromToken)
		if err != nil {
			log.Printf("Error converting user ID: %v", err)
			http.Error(w, "Invalid user ID in token", http.StatusInternalServerError)
			return
		}

		// Perform the deletion using the converted user ID
		if err := data.DeleteUser(client, objectIDFromToken); err != nil {
			log.Printf("Error deleting user: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func extractUserIDFromToken(r *http.Request) (string, error) {
	// Extract the token from the Authorization header
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return "", errors.New("missing Authorization header")
	}

	// Remove 'Bearer ' prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

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
		return "", err
	}

	// Check if the token is valid
	if !token.Valid {
		return "", errors.New("invalid token")
	}

	// Extract user ID from claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("user_id not found in token claims")
	}

	return userID, nil
}

func validateUserInput(user *data.User) error {
	// Validate email format
	if user.Email != "" {
		if user.Email == "aa" {
			return errors.New("invalid email format")
		}
	} else {
		return errors.New("email is required")
	}

	// Other validation logic for other fields

	return nil
}
