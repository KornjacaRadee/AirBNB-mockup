package handlers

import (
	"log"
	"net/smtp"
	"os"
)

func SendEmail(hostEmail, text string) (bool, error) {
	// Sender data
	from := os.Getenv("MAIL_ADDRESS")

	password := os.Getenv("MAIL_PASSWORD")

	// Receiver email
	to := []string{
		hostEmail,
	}

	// smtp server config
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	address := smtpHost + ":" + smtpPort

	// Text
	stringMsg :=
		"AirBNB notification:  \n" +
			"To: " + to[0] + "\n\n" +
			text

	message := []byte(stringMsg)

	// Email Sender Auth
	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(address, auth, from, to, message)
	if err != nil {
		log.Println("Error sending mail", err)
		return false, err
	}
	log.Println("Mail successfully sent")
	return true, nil
}
