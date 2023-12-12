package authHandlers

import (
	"auth-service/client"
	"auth-service/data"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func HandleRegister(dbClient *mongo.Client, pc client.ProfileClient) http.HandlerFunc {
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

		// Check if the password is in the blacklist
		passwordOK, err := data.CheckPasswordInBlacklist(newUser.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !passwordOK {
			http.Error(w, "Password is in the blacklist", http.StatusBadRequest)
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
		if err := data.RegisterUser(dbClient, &newUser); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := pc.SendUserData(newUser); err != nil {
			http.Error(w, "Error sending user data to the profile service", http.StatusInternalServerError)
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
		"roles":   user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func HandleLogin(dbClient *mongo.Client) http.HandlerFunc {
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
		user, err := data.LoginUser(dbClient, credentials.Email, credentials.Password)
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

func HandleGetUserByID(dbClient *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID from URL parameters
		vars := mux.Vars(r)
		userID, err := primitive.ObjectIDFromHex(vars["id"])
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Get user by ID
		user, err := data.GetUserByID(dbClient, userID)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		// Return user data
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func HandleChangePassword(dbClient *mongo.Client) http.HandlerFunc {
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

		user, err := validateOldPasswordAndGetUser(dbClient, userID, passwordChange.OldPassword)
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
		if err := data.UpdatePassword(dbClient, userID, passwordChange.NewPassword); err != nil {
			log.Printf("Error updating password: %v", err)
			http.Error(w, "Error updating password", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// Helper function to validate the old password and get the user
func validateOldPasswordAndGetUser(dbClient *mongo.Client, userID primitive.ObjectID, oldPassword string) (*data.User, error) {
	user, err := data.GetUserByID(dbClient, userID)
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

func HandleGetAllUsers(dbClient *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := data.GetAllUsers(dbClient)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(users)
	}
}

func HandleDeleteUser(dbClient *mongo.Client) http.HandlerFunc {
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
		if err := data.DeleteUser(dbClient, objectIDFromToken); err != nil {
			log.Printf("Error deleting user: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// Request for password recovery
func HandlePasswordRecovery(dbClient *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user's email from the request
		var request struct {
			Email string `json:"email"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}

		// Validate the email and get the user
		user, err := data.GetUserByEmail(dbClient, request.Email)
		if err != nil {
			log.Printf("Error retrieving user: %v", err)
			http.Error(w, "Error retrieving user", http.StatusInternalServerError)
			return
		}

		if user == nil {
			// User not found, but don't disclose this information to the user
			w.WriteHeader(http.StatusOK)
			return
		}

		// Generate a unique recovery token
		recoveryToken, err := data.GenerateRecoveryToken(dbClient, user.ID)
		if err != nil {
			// Handle the error, for example:
			log.Printf("Error generating recovery token: %v", err)
			http.Error(w, "Error generating recovery token", http.StatusInternalServerError)
			return
		}
		// Save the recovery token in the database associated with the user

		// Send a recovery email with a link containing the recovery token
		data.SendRecoveryEmail(user.Email, recoveryToken)

		// Respond to the user indicating that the recovery email has been sent
		w.WriteHeader(http.StatusOK)
	}
}

// Handle the password reset page
func HandlePasswordReset(dbClient *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract recovery token from the URL parameters
		token := r.URL.Query().Get("token")
		if token == "" {
			http.Error(w, "Invalid or missing recovery token", http.StatusBadRequest)
			return
		}

		// Validate the recovery token
		if !data.IsValidRecoveryToken(dbClient, token) {
			http.Error(w, "Invalid recovery token", http.StatusBadRequest)
			return
		}

		// Allow the user to reset their password
		// You can redirect them to a password reset page in your frontend
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]string{"status": "success", "message": "Password reset allowed"}
		json.NewEncoder(w).Encode(response)
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

func HandlePasswordUpdate(dbClient *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract reset token from URL parameters
		resetToken, err := extractResetToken(r)
		if err != nil {
			log.Printf("Error extracting reset token: %v", err)
			http.Error(w, "Invalid reset token", http.StatusBadRequest)
			return
		}

		// Decode the JSON request body
		var passwordChange struct {
			NewPassword string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&passwordChange); err != nil {
			log.Printf("Error decoding JSON: %v", err)
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}
		log.Printf("New password from handler is: %s", passwordChange.NewPassword)

		// Validate the reset token and get the user ID
		userID, err := data.ValidateResetTokenAndGetUser(dbClient, resetToken)
		if err != nil {
			log.Printf("Error validating reset token: %v", err)
			http.Error(w, "Invalid reset token", http.StatusBadRequest)
			return
		}

		// Update the user's password in the database
		if err := data.UpdatePassword(dbClient, userID, passwordChange.NewPassword); err != nil {
			log.Printf("Error updating password: %v", err)
			http.Error(w, "Error updating password", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func validateUserInput(user *data.User) error {
	// Validate email format
	if user.Email != "" {
		if !isValidEmail(user.Email) {
			return errors.New("invalid email format")
		}
	} else {
		return errors.New("email is required")
	}

	// Other validation logic for other fields

	return nil
}

func isValidEmail(email string) bool {
	// Regular expression for basic email validation
	// Note: This regex might not cover all edge cases, consider using a more comprehensive regex if needed
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(emailRegex, email)
	return match
}

func extractResetToken(r *http.Request) (string, error) {
	// Get the token from the "token" query parameter
	token := r.URL.Query().Get("token")

	if token == "" {
		return "", errors.New("reset token not found in the URL")
	}

	return token, nil
}
