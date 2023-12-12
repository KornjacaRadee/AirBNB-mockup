package client

import (
	"auth-service/data"
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type ProfileClient struct {
	address string
	//cb      *gobreaker.CircuitBreaker
}

func NewProfileClient(address string) ProfileClient {
	return ProfileClient{
		address: address,
		//cb:      cb,
	}

}

func (c *ProfileClient) SendUserData(user data.User) error {
	req := convertUser(user)

	reqBytes, err := json.Marshal(req)
	if err != nil {
		log.Println(err)
		return err
	}

	bodyReader := bytes.NewReader(reqBytes)
	requestURL := c.address + "/new"
	httpReq, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)

	if err != nil {
		log.Println(err)
		return errors.New("error sending user data")
	}
	res, err := http.DefaultClient.Do(httpReq)

	if err != nil || res.StatusCode != http.StatusCreated {
		log.Println(err)
		log.Println(res.StatusCode)
		return errors.New("error sending user data")
	}
	return nil
}

func convertUser(user data.User) UserData {
	userData := UserData{
		Name:      user.First_Name,
		Last_Name: user.Last_Name,
		Username:  user.Username,
		Email:     user.Email,
		Address:   user.Address,
		Role:      user.Role,
	}
	return userData
}
