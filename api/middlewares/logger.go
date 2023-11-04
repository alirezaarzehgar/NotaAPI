package middlewares

import (
	"strings"

	"github.com/Asrez/NotaAPI/utils"
	"github.com/labstack/echo/v4"
)

func LogRequests(next echo.HandlerFunc) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		if strings.HasPrefix(c.Request().URL.Path, "/logs") {
			return next(c)
		}
		utils.DebugLog(c.Request().Method, c.Request().URL.Path)
		return next(c)
	})
}
