package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

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
		return utils.ReturnAlert(c, http.StatusInternalServerError, "internal", err)
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

	return c.JSON(http.StatusOK, map[string]any{"status": true, "data": []any{}})
}

func SaveStoryForGuest(c echo.Context) error {
	var story models.Story
	err := db.First(&story, "code", c.Param("code")).Error
	if err == gorm.ErrRecordNotFound {
		return utils.ReturnAlert(c, http.StatusNotFound, "not_found")
	}

	var token models.Token
	err = db.Preload("SavedStories").First(&token, "jwt_token", utils.GetToken(c)).Error
	if err == gorm.ErrRecordNotFound {
		return utils.ReturnAlert(c, http.StatusNotFound, "not_found")
	}

	err = db.Model(&token).Association("SavedStories").Append(&story)
	if err != nil {
		return utils.ReturnAlert(c, http.StatusInternalServerError, "internal")
	}

	return c.JSON(http.StatusOK, map[string]any{"status": true, "data": []any{}})
}

func ListGuestStories(c echo.Context) error {
	conditions := "is_public = ?"
	args := []any{"", true}

	if c.QueryParam("story_type") == models.STORY_TYPE_EXPLORE {
		conditions += " AND `to` IS NULL"
	} else if c.QueryParam("start_date") != "" && c.QueryParam("end_date") != "" {
		startDate, err := time.Parse(time.DateOnly, c.QueryParam("start_date"))
		if err != nil {
			return utils.ReturnAlert(c, http.StatusBadRequest, "bad_request")
		}
		endDate, err := time.Parse(time.DateOnly, c.QueryParam("end_date"))
		if err != nil {
			return utils.ReturnAlert(c, http.StatusBadRequest, "bad_request")
		}

		conditions += " AND `to` >= ? AND `to` <= ?"
		args = append(args, startDate, endDate)
	}

	var token models.Token
	args[0] = conditions
	err := db.Preload("SavedStories", args...).First(&token, "jwt_token", utils.GetToken(c)).Error

	if err == gorm.ErrRecordNotFound {
		return utils.ReturnAlert(c, http.StatusNotFound, "not_found")
	} else if err != nil {
		return utils.ReturnAlert(c, http.StatusNotFound, "internal")
	}

	return c.JSON(http.StatusOK, token.SavedStories)
}

func GuestDeleteAccount(c echo.Context) error {
	var token models.Token

	err := db.First(&token, "jwt_token", utils.GetToken(c)).Error
	if err == gorm.ErrRecordNotFound {
		return utils.ReturnAlert(c, http.StatusNotFound, "not_found")
	} else if err != nil {
		return utils.ReturnAlert(c, http.StatusNotFound, "internal")
	}

	err = db.Model(&token).Association("SavedStories").Clear()
	if err != nil {
		return utils.ReturnAlert(c, http.StatusNotFound, "internal")
	}

	err = db.Delete(&token, token.ID).Error
	log.Println(err)
	if err != nil {
		return utils.ReturnAlert(c, http.StatusNotFound, "internal")
	}

	return c.JSON(http.StatusOK, map[string]any{"status": true, "data": []any{}})
}
