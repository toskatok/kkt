package actions

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// App creates new instance of Echo and configure it
func App(debug bool) *echo.Echo {
	app := echo.New()
	app.Use(middleware.Logger())
	app.Use(middleware.Recover())

	if debug {
		app.Logger.SetLevel(log.DEBUG)
	}

	// Routes
	app.GET("/ping", PingHandler)
	app.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	return app
}
