package client

import (
	"context"
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

func (ac AccommodationClient) CheckIfAccommodationExists(ctx context.Context, id primitive.ObjectID) (bool, error) {

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
		return ac.client.Do(req)
	})
	if err != nil {
		return false, handleHttpReqErr(err, ac.address+"/"+id.Hex(), http.MethodGet, timeout)
	}

	resp := cbResp.(*http.Response)
	if resp.StatusCode != http.StatusOK {
		return false, domain.ErrResp{
			URL:        resp.Request.URL.String(),
			Method:     resp.Request.Method,
			StatusCode: resp.StatusCode,
		}
	}

	return true, nil

	/////////////////////////////////////////////////////////
	//requestURL := client.address + "/" + id.Hex()
	//httpReq, err := http.NewRequest(http.MethodGet, requestURL, nil)
	//
	//if err != nil {
	//	log.Println(err)
	//	return false, err
	//}
	//
	//res, err := http.DefaultClient.Do(httpReq)
	//
	//if err != nil || res.StatusCode != http.StatusOK {
	//	log.Println(err)
	//	log.Println(res.StatusCode)
	//	return false, err
	//}
	//return true, nil
}
