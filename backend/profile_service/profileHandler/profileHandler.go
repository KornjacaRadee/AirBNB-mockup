package profileHandler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/promeneili1/AirBNB-mockup/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"

	"github.com/promeneili1/AirBNB-mockup/data"
	"go.mongodb.org/mongo-driver/mongo"
)

var logger = config.NewLogger("./logging/log.log")

func CreateProfileHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newProfile data.Profile
		if err := json.NewDecoder(r.Body).Decode(&newProfile); err != nil {
			logger.Error("Error decoding request body", map[string]interface{}{"error": err.Error()})
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := data.CreateProfile(client, &newProfile); err != nil {
			logger.Error("Error creating profile", map[string]interface{}{"error": err.Error()})
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Info("Profile created successfully", map[string]interface{}{"profileID": newProfile.ID.Hex()})
		w.WriteHeader(http.StatusCreated)
	}
}

func GetAllProfilesHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profiles, err := data.GetAllProfiles(client)
		if err != nil {
			logger.Error("Error getting all profiles", map[string]interface{}{"error": err.Error()})
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Info("All profiles retrieved successfully", nil)
		json.NewEncoder(w).Encode(profiles)
	}
}

func GetProfileByIDHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		profileID, err := primitive.ObjectIDFromHex(vars["id"])
		if err != nil {
			logger.Error("Invalid profile ID", map[string]interface{}{"error": err.Error()})
			http.Error(w, "Invalid profile ID", http.StatusBadRequest)
			return
		}

		profile, err := data.GetProfileByID(client, profileID)
		if err != nil {
			logger.Error("Error getting profile by ID", map[string]interface{}{"error": err.Error()})
			http.Error(w, "Profile not found", http.StatusNotFound)
			return
		}

		logger.Info("Profile retrieved successfully", map[string]interface{}{"profileID": profileID.Hex()})
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(profile)
	}
}

func UpdateProfileHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		profileID, err := primitive.ObjectIDFromHex(vars["id"])
		if err != nil {
			logger.Error("Invalid profile ID", map[string]interface{}{"error": err.Error()})
			http.Error(w, "Invalid profile ID", http.StatusBadRequest)
			return
		}

		var updatedProfile data.Profile
		if err := json.NewDecoder(r.Body).Decode(&updatedProfile); err != nil {
			logger.Error("Error decoding request body", map[string]interface{}{"error": err.Error()})
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := data.UpdateProfile(client, profileID, &updatedProfile); err != nil {
			logger.Error("Error updating profile", map[string]interface{}{"error": err.Error()})
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Info("Profile updated successfully", map[string]interface{}{"profileID": profileID.Hex()})
		w.WriteHeader(http.StatusOK)
	}
}

func DeleteProfileHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		profileID, err := primitive.ObjectIDFromHex(vars["id"])
		if err != nil {
			logger.Error("Invalid profile ID", map[string]interface{}{"error": err.Error()})
			http.Error(w, "Invalid profile ID", http.StatusBadRequest)
			return
		}

		if err := data.DeleteProfile(client, profileID); err != nil {
			logger.Error("Error deleting profile", map[string]interface{}{"error": err.Error()})
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Info("Profile deleted successfully", map[string]interface{}{"profileID": profileID.Hex()})
		w.WriteHeader(http.StatusOK)
	}
}

func GetProfileByEmailHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		email := vars["email"]
		if email == "" {
			logger.Error("Email parameter missing", map[string]interface{}{"error": "Email parameter missing"})
			http.Error(w, "Email parameter missing", http.StatusBadRequest)
			return
		}

		profile, err := data.GetProfileByEmail(client, email)
		if err != nil {
			logger.Error("Error getting profile by email", map[string]interface{}{"error": err.Error()})
			http.Error(w, "Profile not found", http.StatusNotFound)
			return
		}

		logger.Info("Profile retrieved by email successfully", map[string]interface{}{"email": email})
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(profile)
	}
}
