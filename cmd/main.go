// @title          Macthing Driver API
// @version        1.0
// @description    Finds Nearest Driver According to User Point
// @BasePath       /v1

// @securityDefinitions.apikey InternalAuth
// @in header
// @name Authorization
// @description  Internal calls only.  Format: "Bearer <token>"
package main

import (
	"github.com/mrtuuro/matching-api/internal/application"
	"github.com/mrtuuro/matching-api/internal/client"
	"github.com/mrtuuro/matching-api/internal/config"
	"github.com/mrtuuro/matching-api/internal/router"
	"github.com/mrtuuro/matching-api/internal/service"
)

func main() {
	cfg := config.NewConfig()

	var (
		// CUSTOM HTTP CLIENT INIT
		customClient = client.NewCustomHTTPClient(cfg.DriverAPIUrl, cfg.DriverAPIToken)

		// SERVICE INIT
		matchSvc = service.NewMatchingService(customClient)
	)

	app := application.NewApp(cfg, matchSvc)

	router.Register(app)
	router.PrintRoutes(app)

	app.Run(cfg.Port)
}
