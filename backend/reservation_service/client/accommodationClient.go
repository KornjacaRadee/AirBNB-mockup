package client

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

type AccommodationClient struct {
	address string
}

func NewAccommodationClient(address string) AccommodationClient {
	return AccommodationClient{
		address: address,
	}
}

func (client AccommodationClient) CheckIfAccommodationExists(id primitive.ObjectID) (bool, error) {
	requestURL := client.address + "/" + id.Hex()
	httpReq, err := http.NewRequest(http.MethodGet, requestURL, nil)

	if err != nil {
		log.Println(err)
		return false, err
	}

	res, err := http.DefaultClient.Do(httpReq)

	if err != nil || res.StatusCode != http.StatusOK {
		log.Println(err)
		log.Println(res.StatusCode)
		return false, err
	}
	return true, nil
}
