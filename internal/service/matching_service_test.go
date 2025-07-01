package service

import (
	"context"
	"errors"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/mrtuuro/matching-api/internal/apperror"
	"github.com/mrtuuro/matching-api/internal/client"
	"github.com/mrtuuro/matching-api/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockClient struct{ mock.Mock }

func (m *mockClient) Healthcheck(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}
func (m *mockClient) SearchDriver(ctx context.Context, token string, body echo.Map) (client.SearchResp, error) {
	args := m.Called(ctx, token, body)
	return args.Get(0).(client.SearchResp), args.Error(1)
}

func TestDriverLocationHealthcheck_Error(t *testing.T) {
	ctx := context.Background()
	cl := new(mockClient)
	svc := NewMatchingService(cl)

	clErr := errors.New("upstream down")
	cl.On("Healthcheck", ctx).Return(clErr)

	err := svc.DriverLocationHealthcheck(ctx)

	assert.Same(t, clErr, err)
	cl.AssertExpectations(t)
}

func TestSearchDriver_Success(t *testing.T) {
	ctx := context.Background()
	cl := new(mockClient)
	svc := NewMatchingService(cl)

	want := []model.DriverWithDistance{{
		DriverLocation: model.DriverLocation{DriverID: "d1"},
		DistanceMeters: 123.4,
	}}
	res := client.SearchResp{Data: want}
	body := echo.Map{"dummy": true}

	cl.On("SearchDriver", ctx, "tok", body).Return(res, nil)

	got, err := svc.SearchDriver(ctx, "tok", body)

	assert.NoError(t, err)
	assert.Equal(t, want, got)
	cl.AssertExpectations(t)
}

func TestSearchDriver_NoDriver(t *testing.T) {
	ctx := context.Background()
	cl := new(mockClient)
	svc := NewMatchingService(cl)

	res := client.SearchResp{Data: []model.DriverWithDistance{}}
	cl.On("SearchDriver", ctx, "tok", mock.Anything).Return(res, nil)

	got, err := svc.SearchDriver(ctx, "tok", echo.Map{})

	assert.Nil(t, got)

	appErr, ok := err.(*apperror.AppError)
	if assert.True(t, ok, "error should be AppError") {
		assert.Equal(t, "DRIVER_NEARBY_NOT_FOUND", appErr.Code)
	}
	cl.AssertExpectations(t)
}
