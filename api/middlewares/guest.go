package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/Asrez/NotaAPI/utils"
)

func GuestOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		if utils.GetUserId(c) == 0 {
			return next(c)
		}
		return utils.ReturnAlert(c, http.StatusUnauthorized, "guest_only")
	})
}
