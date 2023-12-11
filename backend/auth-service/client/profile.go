package client

import (
	"auth-service/data"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ProfileClient struct {
	client  *http.Client
	address string
	//cb      *gobreaker.CircuitBreaker
}

func NewProfileClient(client *http.Client, address string) ProfileClient {
	return ProfileClient{
		client:  client,
		address: address,
		//cb:      cb,
	}

}

// Add a function to send user data to the profile service
func (c *ProfileClient) SendUserData(user *data.User) error {
	url := fmt.Sprintf("%s/profiles", c.address) // Adjust the URL according to your profile service routes

	// Convert user data to JSON
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}

	// Make a POST request to the profile service
	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(userJSON))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

//func (pc *ProfileClient) SendUserDataToProfile(user *data.User) error {
//	// Prepare the request body
//	body, err := json.Marshal(user)
//	if err != nil {
//		return err
//	}
//
//	// Create a request to the profile service
//	req, err := http.NewRequest("POST", fmt.Sprintf("%s/profiles", pc.address), bytes.NewBuffer(body))
//	if err != nil {
//		return err
//	}
//	req.Header.Set("Content-Type", "application/json")
//
//	// Send the request
//	resp, err := pc.client.Do(req)
//	if err != nil {
//		return err
//	}
//	defer resp.Body.Close()
//
//	// Check the response status code
//	if resp.StatusCode != http.StatusOK {
//		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
//	}
//
//	return nil
//}
