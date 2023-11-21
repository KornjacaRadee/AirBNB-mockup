package authHandlers

import (
	"auth-service/data"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

//func HandleRegister(client *mongo.Client) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		var newUser data.User
//		if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
//			http.Error(w, err.Error(), http.StatusBadRequest)
//			return
//		}
//
//		// Validate user input
//		if err := validateUserInput(&newUser); err != nil {
//			http.Error(w, err.Error(), http.StatusBadRequest)
//			return
//		}
//
//		hashedPassword, err := hashPassword(newUser.Password)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//		newUser.Password = hashedPassword
//
//		// Register the user
//		if err := data.RegisterUser(client, &newUser); err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//
//		w.WriteHeader(http.StatusCreated)
//	}
//}

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
		hashedPassword, err := hashPassword(newUser.Password)
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expiration time (e.g., 24 hours)
	})

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
		vars := mux.Vars(r)
		userID, err := primitive.ObjectIDFromHex(vars["id"])
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		if err := data.DeleteUser(client, userID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func validateUserInput(user *data.User) error {
	// Implement validation logic here
	return nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
