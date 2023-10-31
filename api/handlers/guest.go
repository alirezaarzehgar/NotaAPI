package handlers

import (
	"encoding/json"
	"io"
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

func SaveStoryForGuest(c echo.Context) error {
	var story models.Story
	r := db.First(&story, "code", c.Param("code"))
	if r.RowsAffected == 0 {
		return utils.ReturnAlert(c, http.StatusNotFound, "not_found")
	}

	err := db.Create(&models.Guest{JwtToken: utils.GetToken(c), StoryCode: story.Code, StoryTo: story.To}).Error
	if err != nil {
		return utils.ReturnAlert(c, http.StatusConflict, "conflict")
	}
	return c.JSON(http.StatusOK, map[string]any{
		"status": true,
		"data":   map[string]any{},
	})
}
func ListGuestStories(c echo.Context) error {
	var guests []models.Guest
	dateCond := db.Where("1 = 1")
	defaultCond := map[string]any{"jwt_token": utils.GetToken(c)}

	if c.QueryParam("story_type") == models.STORY_TYPE_EXPLORE {
		defaultCond["story_to"] = time.Time{}
	} else if c.QueryParam("start_date") != "" && c.QueryParam("end_date") != "" {
		startDate, err := time.Parse(time.DateOnly, c.QueryParam("start_date"))
		if err != nil {
			return utils.ReturnAlert(c, http.StatusBadRequest, "bad_request")
		}
		endDate, err := time.Parse(time.DateOnly, c.QueryParam("end_date"))
		if err != nil {
			return utils.ReturnAlert(c, http.StatusBadRequest, "bad_request")
		}
		dateCond = db.Where("`story_to` >= ? AND `story_to` <= ?", startDate, endDate)
	}

	r := db.Preload("Story").Where(dateCond).Where(defaultCond).Find(&guests)
	if r.Error == gorm.ErrRecordNotFound {
		return utils.ReturnAlert(c, http.StatusNotFound, "not_found")
	} else if r.Error != nil {
		return utils.ReturnAlert(c, http.StatusInternalServerError, "internal")
	}

	var stories []models.Story
	for _, guest := range guests {
		stories = append(stories, guest.Story)
	}

	return c.JSON(http.StatusOK, map[string]any{"status": true, "data": stories})
}
