package client

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/sony/gobreaker"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"reservation_service/domain"
	"time"
)

type AccommodationClient struct {
	client  *http.Client
	address string
	cb      *gobreaker.CircuitBreaker
}

func NewAccommodationClient(client *http.Client, address string, cb *gobreaker.CircuitBreaker) AccommodationClient {
	return AccommodationClient{
		client:  client,
		address: address,
		cb:      cb,
	}
}

func (ac AccommodationClient) GetAccommodation(ctx context.Context, id primitive.ObjectID) (*AccommodationData, error) {
	var timeout time.Duration
	deadline, reqHasDeadline := ctx.Deadline()
	if reqHasDeadline {
		timeout = time.Until(deadline)
	}

	cbResp, err := ac.cb.Execute(func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, ac.address+"/"+id.Hex(), nil)
		if err != nil {
			return nil, err
		}
		resp, err := ac.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, domain.ErrResp{
				URL:        resp.Request.URL.String(),
				Method:     resp.Request.Method,
				StatusCode: resp.StatusCode,
			}
		}

		var accomm AccommodationData

		if err := json.NewDecoder(resp.Body).Decode(&accomm); err != nil {
			return nil, err
		}

		return accomm, nil
	})
	if err != nil {
		return nil, handleHttpReqErr(err, ac.address+"/"+id.Hex(), http.MethodGet, timeout)
	}

	accomm, ok := cbResp.(AccommodationData)
	if !ok {
		return nil, errors.New("invalid response type")
	}

	return &accomm, nil
}
