package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/mrtuuro/matching-api/internal/application"
	"github.com/mrtuuro/matching-api/internal/code"
	"github.com/mrtuuro/matching-api/internal/response"
)

// HealthcheckHandler godoc
// @Summary      Liveness probe
// @Description  Returns 200 OK with a success 
// @Tags         system
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.SwaggerSuccess
// @Router       /v1/healthz [get]
func HealthcheckHandler(app *application.Application) echo.HandlerFunc {
	return func(c echo.Context) error {
		return response.RespondSuccess[any](c, code.SuccessHealthCheck, nil)
	}
}
