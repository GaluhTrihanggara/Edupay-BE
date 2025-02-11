package middlewares

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func LogMiddlewares(e *echo.Echo) {
	var (
		time    = "${time_rfc3339}"
		status  = "${status}"
		method  = "${method}"
		host    = "${host}"
		path    = "${path}"
		latency = "${latency_human}"
	)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: fmt.Sprintf(`"time":"%s" %s "status":%s" %s "method":"%s" %s "host":"%s" %s "path":"%s" %s "latency_human":"%s"	%s%s`, time, "\n", status, "\n", method, "\n", host, "\n", path, "\n", latency, "\n", "\n"),
	}))
}
