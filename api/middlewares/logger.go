package middlewares

import (
	"encoding/json"
	"strings"

	"github.com/Asrez/NotaAPI/utils"
	"github.com/labstack/echo/v4"
)

func LogRequests(next echo.HandlerFunc) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		if strings.HasPrefix(c.Request().URL.Path, "/logs") {
			return next(c)
		}

		var data map[string]any

		var output []any
		output = append(output, c.Request().Method, c.Request().URL.Path)

		if data != nil {
			output = append(output, ":", data)
		}

		json.NewDecoder(c.Request().Body).Decode(&data)
		utils.DebugLog(output...)
		return next(c)
	})
}
