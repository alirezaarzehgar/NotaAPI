package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/Asrez/NotaAPI/models"
	"github.com/Asrez/NotaAPI/utils"
)

func CheckToken(c echo.Context) error {
	token := models.Token{}
	if err := json.NewDecoder(c.Request().Body).Decode(&token); err != nil {
		return utils.ReturnAlert(c, http.StatusBadRequest, "bad_request")
	}

	searchPattern := map[string]any{"jwt_token": token.JwtToken, "blocked": false}
	if db.Where(searchPattern).First(&token).RowsAffected == 0 {
		return utils.ReturnAlert(c, http.StatusNotFound, "token_not_found")
	}

	return c.JSON(http.StatusOK, map[string]any{"status": true})
}
