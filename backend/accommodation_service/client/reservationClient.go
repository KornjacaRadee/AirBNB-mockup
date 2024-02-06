package client

import (
	"accommodation_service/domain"
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

func (rc ReservationClient) CheckDates(ctx context.Context, req SearchReqs) ([]*primitive.ObjectID, error) {
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	var timeout time.Duration
	deadline, reqHasDeadline := ctx.Deadline()
	if reqHasDeadline {
		timeout = time.Until(deadline)
	}

	cbResp, err := rc.cb.Execute(func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, rc.address+"/check-dates", bytes.NewBuffer(reqBytes))
		if err != nil {
			return nil, err
		}
		resp, err := rc.client.Do(req)
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

		var ids []*primitive.ObjectID
		if err := json.NewDecoder(resp.Body).Decode(&ids); err != nil {
			return nil, err
		}

		return ids, nil
	})

	if err != nil {
		return nil, handleHttpReqErr(err, rc.address+"/check-dates", http.MethodGet, timeout)
	}

	ids, ok := cbResp.([]*primitive.ObjectID)
	if !ok {
		return nil, errors.New("invalid response type")
	}

	return ids, nil
}
