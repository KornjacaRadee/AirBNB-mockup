package profileHandler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"

	"github.com/promeneili1/AirBNB-mockup/data"

	"go.mongodb.org/mongo-driver/mongo"
)

func CreateProfileHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newProfile data.Profile
		if err := json.NewDecoder(r.Body).Decode(&newProfile); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := data.CreateProfile(client, &newProfile); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func GetAllProfilesHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profiles, err := data.GetAllProfiles(client)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(profiles)
	}
}

func GetProfileByIDHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract profile ID from URL parameters
		vars := mux.Vars(r)
		profileID, err := primitive.ObjectIDFromHex(vars["id"])
		if err != nil {
			http.Error(w, "Invalid profile ID", http.StatusBadRequest)
			return
		}

		// Get profile by ID
		profile, err := data.GetProfileByID(client, profileID)
		if err != nil {
			http.Error(w, "Profile not found", http.StatusNotFound)
			return
		}

		// Return profile data
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(profile)
	}
}

func UpdateProfileHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract profile ID from URL parameters
		vars := mux.Vars(r)
		profileID, err := primitive.ObjectIDFromHex(vars["id"])
		if err != nil {
			http.Error(w, "Invalid profile ID", http.StatusBadRequest)
			return
		}

		// Decode the JSON request body
		var updatedProfile data.Profile
		if err := json.NewDecoder(r.Body).Decode(&updatedProfile); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Update the profile
		if err := data.UpdateProfile(client, profileID, &updatedProfile); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func DeleteProfileHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract profile ID from URL parameters
		vars := mux.Vars(r)
		profileID, err := primitive.ObjectIDFromHex(vars["id"])
		if err != nil {
			http.Error(w, "Invalid profile ID", http.StatusBadRequest)
			return
		}

		// Delete the profile
		if err := data.DeleteProfile(client, profileID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func GetProfileByEmailHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract email from URL parameters or request body
		vars := mux.Vars(r)
		var email string
		email = vars["email"]
		if email == "" {
			http.Error(w, "Email parameter missing", http.StatusBadRequest)
			return
		}

		// Get profile by email
		profile, err := data.GetProfileByEmail(client, email)
		if err != nil {
			http.Error(w, "Profile not found", http.StatusNotFound)
			return
		}

		// Return profile data
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(profile)
	}
}
