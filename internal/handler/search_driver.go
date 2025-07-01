package handler

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/mrtuuro/matching-api/internal/apperror"
	"github.com/mrtuuro/matching-api/internal/application"
	"github.com/mrtuuro/matching-api/internal/code"
	"github.com/mrtuuro/matching-api/internal/model"
	"github.com/mrtuuro/matching-api/internal/response"
)

type SearchDriverRequest struct {
	Location GeoPointDTO `json:"location" validate:"required"`
	Radius   float64     `json:"radius" validate:"required,gt=0"`
	Limit    int         `json:"limit" validate:"required,gte=1"`
}

type GeoPointDTO struct {
	Type        string     `json:"type" validate:"required,eq=Point"`
	Coordinates [2]float64 `json:"coordinates" validate:"required,lte=180,gte=-180,dive"`
}

// SearchDriverHandler godoc
// @Summary      Find nearest drivers
// @Description  Returns drivers ordered by distance; distance (metres) is pre-calculated.
// @Tags         drivers
// @Accept       json
// @Produce      json
// @Param        body  body  handler.SearchDriverRequest  true  "Search parameters"
// @Success      200   {object}  response.SwaggerSuccess        "List of DriverWithDistance"
// @Failure      400   {object}  response.SwaggerError
// @Failure      401   {object}  response.SwaggerError
// @Failure      404   {object}  response.SwaggerError
// @Failure      500   {object}  response.SwaggerError
// @Security     InternalAuth
// @Router       /v1/drivers/search [post]
func SearchDriverHandler(app *application.Application) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := context.WithoutCancel(c.Request().Context())
		var req SearchDriverRequest
		if err := c.Bind(&req); err != nil {
			return response.RespondError[any](c, apperror.NewAppError(
				code.ErrInvalidJSON,
				err,
				code.GetErrorMessage(code.ErrInvalidJSON),
				))
		}

		if err := c.Validate(&req); err != nil {
			return response.RespondError[any](c, apperror.NewAppError(
				code.ErrValidationFailed,
				err,
				code.GetErrorMessage(code.ErrValidationFailed),
				))
		}

		authToken := c.Get("authToken")
		body := echo.Map{
			"location": req.Location,
			"radius":   req.Radius,
			"limit":    req.Limit,
		}

		drivers, err := app.MatchingService.SearchDriver(ctx, authToken.(string), body)
		if err != nil {
			return response.RespondError[any](c, err)
		}

		return response.RespondSuccess[[]model.DriverWithDistance](c, code.SuccessOperationCompleted, &drivers)

	}
}
