package client

import (
	"accommodation_service/domain"
	"context"
	"github.com/sony/gobreaker"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
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

func (ac ReservationClient) CheckReservationsByAccommodationId(ctx context.Context, id primitive.ObjectID) (bool, error) {
	var timeout time.Duration
	deadline, reqHasDeadline := ctx.Deadline()
	if reqHasDeadline {
		timeout = time.Until(deadline)
	}

	cbResp, err := ac.cb.Execute(func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, ac.address+"/accomm/"+id.Hex()+"/check", nil)
		if err != nil {
			return true, err
		}
		return ac.client.Do(req)
	})
	if err != nil {
		return true, handleHttpReqErr(err, ac.address+"/accomm/"+id.Hex()+"/check", http.MethodGet, timeout)
	}

	resp := cbResp.(*http.Response)
	if resp.StatusCode != http.StatusOK {
		return true, domain.ErrResp{
			URL:        resp.Request.URL.String(),
			Method:     resp.Request.Method,
			StatusCode: resp.StatusCode,
		}
	}

	return false, nil
}
