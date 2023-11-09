package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/Asrez/NotaAPI/models"
	"github.com/Asrez/NotaAPI/utils"
)

func AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		var adminCnt int64
		db.First(models.User{}, "id = ? AND role = ?", utils.GetUserId(c), models.USERS_ROLE_ADMIN).Count(&adminCnt)

		if adminCnt >= 0 {
			return next(c)
		}
		return utils.ReturnAlert(c, http.StatusUnauthorized, "admin_only")

	})
}
