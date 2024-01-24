package client

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/sony/gobreaker"
	"net/http"
	"notification_service/domain"
	"time"
)

type ProfileClient struct {
	client  *http.Client
	address string
	cb      *gobreaker.CircuitBreaker
}

func NewProfileClient(client *http.Client, address string, cb *gobreaker.CircuitBreaker) ProfileClient {
	return ProfileClient{
		client:  client,
		address: address,
		cb:      cb,
	}
}

func (uc ProfileClient) GetAllInformationsByUserID(ctx context.Context, id string) (*UserData, error) {
	var timeout time.Duration
	deadline, reqHasDeadline := ctx.Deadline()
	if reqHasDeadline {
		timeout = time.Until(deadline)
	}

	cbResp, err := uc.cb.Execute(func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, uc.address+"/"+id, nil)
		if err != nil {
			return nil, err
		}
		resp, err := uc.client.Do(req)
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

		var user UserData
		if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
			return nil, err
		}

		return user, nil
	})

	if err != nil {
		return nil, handleHttpReqErr(err, uc.address+"/"+id, http.MethodGet, timeout)
	}

	user, ok := cbResp.(UserData)
	if !ok {
		return nil, errors.New("invalid response type")
	}

	return &user, nil
}
