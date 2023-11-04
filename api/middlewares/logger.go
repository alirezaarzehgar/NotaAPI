package middlewares

import (
	"encoding/json"

	"github.com/labstack/echo/v4"

	"github.com/Asrez/NotaAPI/utils"
)

func LogRequests(next echo.HandlerFunc) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		var data map[string]any
		json.NewDecoder(c.Request().Body).Decode(&data)
		utils.DebugLog(c.Request().URL.Path, ":", data)
		return next(c)
	})
}
