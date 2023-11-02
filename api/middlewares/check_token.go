package middlewares

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/Asrez/NotaAPI/models"
	"github.com/Asrez/NotaAPI/utils"
)

func CheckToken(next echo.HandlerFunc) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		var token models.Token

		r := db.Where(models.Token{JwtToken: utils.GetToken(c)}).First(&token)
		if r.RowsAffected == 0 || token.Blocked {
			return utils.ReturnAlert(c, http.StatusUnauthorized, "invalid_token")
		}

		err := db.Where(token).Updates(models.Token{LastRequestTime: time.Now()}).Error
		if err != nil {
			return utils.ReturnAlert(c, http.StatusInternalServerError, "internal")
		}

		return next(c)
	})
}
