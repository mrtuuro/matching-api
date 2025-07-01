package handler

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mrtuuro/matching-api/internal/application"
	"github.com/mrtuuro/matching-api/internal/code"
	"github.com/mrtuuro/matching-api/internal/response"
)

// HealthcheckHandler godoc
// @Summary      Liveness probe
// @Description  Checks if Driver Location API is up and running. Returns 200 OK with a success 
// @Tags         system
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.SwaggerSuccess
// @Router       /v1/driver-healthcheck [get]
func DriverLocationHealthcheckHandler(app *application.Application) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
		defer cancel()
		if err := app.MatchingService.DriverLocationHealthcheck(ctx); err != nil {
			return response.RespondError[any](c, err)
		}
		return response.RespondSuccess[any](c, code.SuccessHealthCheck, nil)
	}
}
