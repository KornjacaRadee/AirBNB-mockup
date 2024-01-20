package clients

import (
	"context"
	"github.com/sony/gobreaker"
	"log"
	"net/http"
)

type AccommodationsClient struct {
	client  *http.Client
	address string
	cb      *gobreaker.CircuitBreaker
}

func NewAccommodationsClient(client *http.Client, address string, cb *gobreaker.CircuitBreaker) AccommodationsClient {
	return AccommodationsClient{
		client:  client,
		address: address,
		cb:      cb,
	}

}
func (ac AccommodationsClient) DeleteUserAccommodations(ctx context.Context, id string) (interface{}, error) {
	cbResp, err := ac.cb.Execute(func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodDelete, ac.address+"/profiles/"+id, http.NoBody)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		return ac.client.Do(req)
	})
	if err != nil {
		return errors.NewError("internal error", 500)
	}
	resp := cbResp.(*http.Response)
	if resp.StatusCode == 201 {
		return nil
	}
	return errors.NewError("internal error", resp.StatusCode)
}
