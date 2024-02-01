package client

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/sony/gobreaker"
	"log"
	"net/http"
	"rating_service/domain"
	"time"
)

type NotificationClient struct {
	client  *http.Client
	address string
	cb      *gobreaker.CircuitBreaker
}

func NewNotificationClient(client *http.Client, address string, cb *gobreaker.CircuitBreaker) NotificationClient {
	return NotificationClient{
		client:  client,
		address: address,
		cb:      cb,
	}
}

func (c *NotificationClient) SendReservationNotification(ctx context.Context, notification NotificationData) (interface{}, error) {
	reqBytes, err := json.Marshal(notification)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var timeout time.Duration
	deadline, reqHasDeadline := ctx.Deadline()
	if reqHasDeadline {
		timeout = time.Until(deadline)
	}

	cbResp, err := c.cb.Execute(func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.address+"/new", bytes.NewBuffer(reqBytes))
		if err != nil {
			return nil, err
		}
		return c.client.Do(req)
	})
	if err != nil {
		return nil, handleHttpReqErr(err, c.address, http.MethodPost, timeout)
	}
	resp := cbResp.(*http.Response)
	if resp.StatusCode != http.StatusCreated {
		return nil, domain.ErrResp{
			URL:        resp.Request.URL.String(),
			Method:     resp.Request.Method,
			StatusCode: resp.StatusCode,
		}
	}

	return true, nil
}
