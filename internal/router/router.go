package router

import (
	"fmt"

	"github.com/mrtuuro/matching-api/internal/application"
	"github.com/mrtuuro/matching-api/internal/handler"
	"github.com/mrtuuro/matching-api/internal/middleware"
)

func Register(app *application.Application) {
	app.E.GET("/healthz", handler.HealthcheckHandler(app)).Name = "Healthcheck"

	protected := app.E.Group("/v1")
	protected.Use(middleware.CustomMiddleware(app))

	protected.GET("/driver-healthcheck", handler.DriverLocationHealthcheckHandler(app)).Name = "Driver Location API Healthcheck"
	protected.POST("/drivers/search", handler.SearchDriverHandler(app)).Name = "Search Driver With Given Location"

}

func PrintRoutes(app *application.Application) {
	fmt.Println("=== ROUTES ===")
	routes := app.E.Routes()
	for _, r := range routes {
		fmt.Printf("%s - [%s]%s\n", r.Name, r.Method, r.Path)
	}
}
