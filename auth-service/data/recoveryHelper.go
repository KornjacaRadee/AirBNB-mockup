package data

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/mail"
	"net/smtp"
	"strconv"
)

func generateUniqueToken() string {
	// Generate a random byte slice (e.g., 32 bytes)
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		// Handle error (e.g., log, return an error)
		return ""
	}

	// Encode the random byte slice to base64 to create a string token
	token := base64.URLEncoding.EncodeToString(randomBytes)

	return token
}

func GenerateRecoveryToken(client *mongo.Client, userID primitive.ObjectID) (string, error) {
	resetToken := generateUniqueToken()

	// Store the reset token along with the user ID in a database collection
	resetTokenCollection := client.Database("authDB").Collection("reset_tokens")
	tokenDoc := bson.D{
		{"token", resetToken},
		{"user_id", userID},
	}
	_, err := resetTokenCollection.InsertOne(context.TODO(), tokenDoc)
	if err != nil {
		return "", err
	}

	return resetToken, nil
}

func SendRecoveryEmail(toEmail, recoveryToken string) {
	// Set up your SMTP configuration
	smtpServer := "smtp.ethereal.email"
	smtpPort := 587
	smtpUsername := "caesar.graham6@ethereal.email"
	smtpPassword := "d3jsKjjcQrfbzXyMqs"

	// Set up the email message
	from := mail.Address{"no-reply", "passrecovery@airbnbb.com"}
	to := mail.Address{"", toEmail}
	subject := "Password Recovery"
	body := "Click the following link to reset your password: http://localhost:8082/reset?token=" + recoveryToken

	// Connect to the SMTP server
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpServer)
	_, err := smtp.Dial(smtpServer + ":" + strconv.Itoa(smtpPort))
	if err != nil {
		log.Printf("Error connecting to SMTP server: %v", err)
		return
	}

	// Set up the email headers
	message := []byte("To: " + to.String() + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\r\n\r\n" +
		body)

	// Send the email
	if err := smtp.SendMail(smtpServer+":"+strconv.Itoa(smtpPort), auth, from.Address, []string{to.Address}, message); err != nil {
		log.Printf("Error sending recovery email: %v", err)
		return
	}
}
func IsValidRecoveryToken(client *mongo.Client, resetToken string) bool {
	// Query the database to find the user ID associated with the reset token
	resetTokenCollection := client.Database("authDB").Collection("reset_tokens")
	var tokenDoc bson.M
	err := resetTokenCollection.FindOne(context.TODO(), bson.D{{"token", resetToken}}).Decode(&tokenDoc)
	if err != nil {
		// Token not found or other error, consider it invalid
		return false
	}

	// Extract the user ID from the tokenDoc
	_, ok := tokenDoc["user_id"].(primitive.ObjectID)
	if !ok {
		// User ID not found in reset token document, consider it invalid
		return false
	}

	// Optional: Add additional checks, such as token expiration, if needed

	return true
}

// Validate the recovery token format
func ValidateResetTokenAndGetUser(client *mongo.Client, resetToken string) (primitive.ObjectID, error) {
	// Query the database to find the user ID associated with the reset token
	resetTokenCollection := client.Database("authDB").Collection("reset_tokens")
	var tokenDoc bson.M
	err := resetTokenCollection.FindOne(context.TODO(), bson.D{{"token", resetToken}}).Decode(&tokenDoc)
	if err != nil {
		log.Printf("Error querying reset token: %v", err)
		return primitive.NilObjectID, err
	}

	// Log the retrieved token document
	log.Printf("Retrieved Token Document: %+v", tokenDoc)

	// Extract the user ID from the tokenDoc
	userID, ok := tokenDoc["user_id"].(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, errors.New("user ID not found in reset token document")
	}

	// Optional: Delete the reset token from the database after it's used
	_, err = resetTokenCollection.DeleteOne(context.TODO(), bson.D{{"token", resetToken}})
	if err != nil {
		log.Printf("Error deleting reset token: %v", err)
	}

	return userID, nil
}
