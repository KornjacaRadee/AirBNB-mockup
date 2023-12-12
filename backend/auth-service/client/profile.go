package client

import (
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

func (c *ProfileClient) SendUserData(userData UserData) error {
	req := userData

	reqBytes, err := json.Marshal(req)
	if err != nil {
		log.Println(err)
		return err
	}

	bodyReader := bytes.NewReader(reqBytes)
	requestURL := c.address + "/register"
	httpReq, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)

	if err != nil {
		log.Println(err)
		return errors.New("error sending user data")
	}
	res, err := http.DefaultClient.Do(httpReq)

	if err != nil || res.StatusCode != http.StatusOK {
		log.Println(err)
		log.Println(res.StatusCode)
		return errors.New("error sending user data")
	}
	return nil
}
