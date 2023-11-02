package middlewares

import (
	"net/http"

	"github.com/Asrez/NotaAPI/utils"
	"github.com/labstack/echo/v4"
)

func UserOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		if utils.GetUserId(c) > 0 {
			return next(c)
		}
		return utils.ReturnAlert(c, http.StatusUnauthorized, "user_only")
	})
}
