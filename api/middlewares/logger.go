package middlewares

import (
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/Asrez/NotaAPI/utils"
)

func LogRequests(next echo.HandlerFunc) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		if !strings.HasPrefix(c.Request().URL.Path, "/logs") {
			utils.DebugLog(c.Request().Method, c.Request().URL.Path)
		}
		return next(c)
	})
}
