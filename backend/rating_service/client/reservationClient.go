package client

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/sony/gobreaker"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"rating_service/domain"
	"time"
)

type ReservationClient struct {
	client  *http.Client
	address string
	cb      *gobreaker.CircuitBreaker
}

func NewReservationClient(client *http.Client, address string, cb *gobreaker.CircuitBreaker) ReservationClient {
	return ReservationClient{
		client:  client,
		address: address,
		cb:      cb,
	}
}

func (ac ReservationClient) GetReservationsByGuestId(ctx context.Context, id primitive.ObjectID) (ReservationsData, error) {
	var timeout time.Duration
	deadline, reqHasDeadline := ctx.Deadline()
	if reqHasDeadline {
		timeout = time.Until(deadline)
	}

	cbResp, err := ac.cb.Execute(func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, ac.address+"/guest/"+id.Hex()+"/reservations", nil)
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

		var reservations ReservationsData
		if err := json.NewDecoder(resp.Body).Decode(&reservations); err != nil {
			return nil, err
		}

		return reservations, nil
	})
	if err != nil {
		return nil, handleHttpReqErr(err, ac.address+"/guest/"+id.Hex()+"/reservations", http.MethodGet, timeout)
	}

	reservations, ok := cbResp.(ReservationsData)
	if !ok {
		return nil, errors.New("invalid response type")
	}

	return reservations, nil
}
