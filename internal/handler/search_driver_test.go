package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/mrtuuro/matching-api/internal/application"
	"github.com/mrtuuro/matching-api/internal/model"
	"github.com/mrtuuro/matching-api/internal/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type dummyValidator struct{}

func (dummyValidator) Validate(i interface{}) error { return nil }

type stubMatchSvc struct{}

func (stubMatchSvc) Healthcheck(ctx context.Context) error {
	return nil
}
func (stubMatchSvc) DriverLocationHealthcheck(ctx context.Context) error {
	return nil
}

func (stubMatchSvc) SearchDriver(
	ctx context.Context,
	token string,
	body echo.Map,
) ([]model.DriverWithDistance, error) {

	return []model.DriverWithDistance{{
		DriverLocation: model.DriverLocation{DriverID: "d1"},
		DistanceMeters: 42.0,
	}}, nil
}

func TestSearchDriverHandler_HappyPath(t *testing.T) {
	e := echo.New()
	e.Validator = dummyValidator{}
	app := &application.Application{MatchingService: stubMatchSvc{}}

	h := SearchDriverHandler(app)

	reqJSON := `{
	"location": { "type": "Point", "coordinates": [10,10] },
	"radius":   1000,
	"limit":    1
	}`

	req := httptest.NewRequest(http.MethodPost, "/v1/drivers/search", bytes.NewBufferString(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	// middleware stored authToken
	ctx.Set("authToken", "user-jwt")

	_ = h(ctx)

	assert.Equal(t, http.StatusOK, rec.Code)

	var got response.APIResponse[[]model.DriverWithDistance]
	_ = json.Unmarshal(rec.Body.Bytes(), &got)

	assert.True(t, got.Success)
	require.NotNil(t, got.Data)
	drivers := *got.Data

	require.Len(t, drivers, 1)
	assert.Equal(t, "d1", drivers[0].DriverID)
	assert.Equal(t, 42.0, drivers[0].DistanceMeters)
}
