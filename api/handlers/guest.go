package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/Asrez/NotaAPI/models"
	"github.com/Asrez/NotaAPI/utils"
)

func CreateGuestToken(c echo.Context) error {
	var tokenMap map[string]any
	var token models.Token
	body, _ := io.ReadAll(c.Request().Body)

	if err := json.Unmarshal(body, &tokenMap); err != nil {
		return utils.ReturnAlert(c, http.StatusBadRequest, "bad_request")
	}

	for _, field := range []string{"screen_height", "screen_width", "resolution", "device_type", "version"} {
		if tokenMap[field] == nil {
			return utils.ReturnAlert(c, http.StatusBadRequest, "bad_request")
		}
	}

	json.Unmarshal(body, &token)
	token.JwtToken = utils.CreateGuestToken()
	if err := db.Create(&token).Error; err != nil {
		return utils.ReturnAlert(c, http.StatusInternalServerError, "internal")
	}

	return c.JSON(http.StatusOK, map[string]any{
		"status": true,
		"data":   map[string]any{"token": token.JwtToken},
	})
}

func GetGuestSettings(c echo.Context) error {
	var token models.Token
	r := db.Model(&models.Token{}).
		Where(models.Token{JwtToken: utils.GetToken(c)}).
		First(&token)
	if r.RowsAffected == 0 {
		return utils.ReturnAlert(c, http.StatusNotFound, "not_found")
	}
	return c.JSON(http.StatusOK, map[string]any{
		"status": true,
		"data": map[string]any{
			"notification": token.Notification,
			"gcm_token":    token.GCMToken,
		},
	})
}

func EditGuestSettings(c echo.Context) error {
	var token models.Token

	if err := json.NewDecoder(c.Request().Body).Decode(&token); err != nil {
		return utils.ReturnAlert(c, http.StatusBadRequest, "bad_request")
	}

	r := db.Model(&models.Token{}).
		Where(models.Token{JwtToken: utils.GetToken(c)}).
		Updates(map[string]any{"notification": token.Notification, "gcm_token": token.GCMToken})
	if r.RowsAffected == 0 {
		return utils.ReturnAlert(c, http.StatusNotFound, "not_found")
	}

	return c.JSON(http.StatusOK, map[string]any{
		"status": true,
		"data":   map[string]any{},
	})
}
