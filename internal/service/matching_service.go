package service

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/mrtuuro/matching-api/internal/apperror"
	"github.com/mrtuuro/matching-api/internal/client"
	"github.com/mrtuuro/matching-api/internal/model"
)

type MatchingService interface {
	Healthcheck(ctx context.Context) error
	DriverLocationHealthcheck(ctx context.Context) error
	SearchDriver(ctx context.Context, token string, body echo.Map) ([]model.DriverWithDistance, error)
}

type matchingService struct {
	client client.CustomClient
}

func NewMatchingService(cl client.CustomClient) MatchingService {
	svc := &matchingService{client: cl}
	return svc
}

func (s *matchingService) Healthcheck(ctx context.Context) error {
	return nil
}

func (s *matchingService) DriverLocationHealthcheck(ctx context.Context) error {
	return s.client.Healthcheck(ctx)
}

func (s *matchingService) SearchDriver(ctx context.Context, token string, body echo.Map) ([]model.DriverWithDistance, error) {
	resp, err := s.client.SearchDriver(ctx, token, body)
	if err != nil {
		return nil, err
	}

	drivers := resp.Data
	if len(drivers) == 0 {
		return nil, apperror.NewAppError(
			"DRIVER_NEARBY_NOT_FOUND",
			nil,
			"There is no available driver found nearby.",
			)
	}

	// drivers, err := s.client.SearchDriver(ctx, token, body)
	// if err != nil {
	// 	return nil, err
	// }
	return drivers, nil
}
