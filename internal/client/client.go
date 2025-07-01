package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mrtuuro/matching-api/internal/apperror"
	"github.com/mrtuuro/matching-api/internal/code"
	"github.com/sony/gobreaker"
)

type CustomClient interface {
	Healthcheck(ctx context.Context) error
	SearchDriver(ctx context.Context, token string, body echo.Map) (dlSearchResp, error)
}

type customHTTPClient struct {
	client *http.Client

	baseURL string
	token   string
	cb      *gobreaker.CircuitBreaker
}

func NewCustomHTTPClient(baseURL, token string) CustomClient {
	cl := &customHTTPClient{
		client:  &http.Client{Timeout: 3 * time.Second},
		baseURL: baseURL,
		token:   token,
	}

	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "Driver Location API",
		Timeout:     5 * time.Second,
		MaxRequests: 3,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.Requests >= 10
		},
	})

	cl.cb = cb
	return cl
}

func (cl *customHTTPClient) Healthcheck(ctx context.Context) error {
	fmt.Println("we are in healtheck request")

	_, err := cl.cb.Execute(func() (any, error) {
		req, _ := http.NewRequestWithContext(ctx, "GET", cl.baseURL+"/healthz", nil)
		req.Header.Set("Authorization", "Bearer "+cl.token)

		resp, err := cl.client.Do(req)
		if err != nil {
			fmt.Println("system internal")
			return nil, apperror.NewAppError(
				code.ErrSystemInternal,
				err,
				code.GetErrorMessage(code.ErrSystemInternal),
				)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, apperror.NewAppError(
				"ERR_DRIVER_API_NOT_RESPONDING",
				fmt.Errorf("Driver location api is unhealthy %v", resp.StatusCode),
				"Driver location api is unhealthy",
				)
		}
		return nil, nil
	})

	return err
}
