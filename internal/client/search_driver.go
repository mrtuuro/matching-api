package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mrtuuro/matching-api/internal/apperror"
	"github.com/mrtuuro/matching-api/internal/code"
	"github.com/mrtuuro/matching-api/internal/model"
)

type SearchResp struct {
	Success bool                       `json:"success"`
	Code    string                     `json:"code"`
	Message string                     `json:"message"`
	Data    []model.DriverWithDistance `json:"data"`
}

func (c *customHTTPClient) SearchDriver(ctx context.Context, token string, body echo.Map) (SearchResp, error) {
	var payload SearchResp

	_, err := c.cb.Execute(func() (any, error) {
		buf, _ := json.Marshal(body)

		req, _ := http.NewRequestWithContext(ctx, http.MethodPost,
			c.baseURL+"/v1/drivers/search", bytes.NewReader(buf))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.client.Do(req)
		if err != nil {
			return nil, apperror.NewAppError(
				code.ErrSystemInternal,
				err,
				code.GetErrorMessage(code.ErrSystemInternal),
				)
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusInternalServerError:
			return nil, apperror.NewAppError(
				code.ErrSystemInternal,
				err,
				code.GetErrorMessage(code.ErrSystemInternal),
				)
		case http.StatusOK:

			if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
				fmt.Println("here")

				return nil, apperror.NewAppError(
					code.ErrSystemInternal,
					err,
					code.GetErrorMessage(code.ErrSystemInternal),
					)
			}
			return nil, nil
		case http.StatusNotFound:
			return []model.DriverWithDistance{}, nil
		default:
			return nil, apperror.NewAppError(
				"STATUS_CODE_UNKNOWN",
				fmt.Errorf("API responded with an unknown status code"),
				"API responded with an unknown status code",
				)
		}
	})
	return payload, err

}
