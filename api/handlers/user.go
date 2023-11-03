package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/Asrez/NotaAPI/models"
	"github.com/Asrez/NotaAPI/utils"
)

func Register(c echo.Context) error {
	user := models.User{Role: models.USERS_ROLE_USER}
	if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
		return utils.ReturnAlert(c, http.StatusBadRequest, "bad_request")
	}
	if user.Email == "" || user.Password == "" {
		return utils.ReturnAlert(c, http.StatusBadRequest, "bad_request")
	}

	if utils.ValidatePassword(user.Password) {
		return utils.ReturnAlert(c, http.StatusBadRequest, "insecure_password")
	}
	if _, err := mail.ParseAddress(user.Email); err != nil {
		return utils.ReturnAlert(c, http.StatusBadRequest, "invalid_email")
	}

	user.Password = utils.CreateSHA256(user.Password)
	err := db.Create(&user).Error
	if err == gorm.ErrDuplicatedKey {
		return utils.ReturnAlert(c, http.StatusConflict, "user_conflict")
	} else if err != nil {
		return utils.ReturnAlert(c, http.StatusInternalServerError, "internal")
	}

	token := utils.CreateUserToken(user.ID, user.Email, user.Username)
	if err := db.Create(&models.Token{UserID: user.ID, JwtToken: token}).Error; err != nil {
		return utils.ReturnAlert(c, http.StatusInternalServerError, "internal")
	}

	return c.JSON(http.StatusOK, map[string]any{
		"status": true,
		"data":   map[string]any{"token": token},
	})
}

func Login(c echo.Context) error {
	var loggedin int64
	var user models.User

	if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
		return utils.ReturnAlert(c, http.StatusBadRequest, "bad_request")
	}
	fillteredUser := models.User{Email: user.Email, Password: utils.CreateSHA256(user.Password)}
	db.Where(fillteredUser).First(&user).Count(&loggedin)

	if loggedin == 0 {
		return utils.ReturnAlert(c, http.StatusUnauthorized, "login_unauthorized")
	}

	token := utils.CreateUserToken(user.ID, user.Email, user.Username)
	if err := db.Create(&models.Token{UserID: user.ID, JwtToken: token}).Error; err != nil {
		return utils.ReturnAlert(c, http.StatusInternalServerError, "internal")
	}

	return c.JSON(http.StatusOK, map[string]any{
		"status": true,
		"data":   map[string]any{"token": token},
	})
}

func GetStoryCount(c echo.Context) error {
	var storyNormalCount, storyExploreCount int64

	err := db.Model(&models.Story{}).
		Where(models.Story{UserID: utils.GetUserId(c), Type: models.STORY_TYPE_NORMAL}).
		Count(&storyNormalCount).Error
	if err != nil {
		return utils.ReturnAlert(c, http.StatusInternalServerError, "internal")
	}

	err = db.Model(&models.Story{}).
		Where(models.Story{UserID: utils.GetUserId(c), Type: models.STORY_TYPE_EXPLORE}).
		Count(&storyExploreCount).Error
	if err != nil {
		return utils.ReturnAlert(c, http.StatusInternalServerError, "internal")
	}

	return c.JSON(http.StatusOK, map[string]any{
		"status": true,
		"data": map[string]any{
			"explore_story":      storyExploreCount > 0,
			"normal_story_count": storyNormalCount,
		},
	})
}

func UserDeleteAccount(c echo.Context) error {
	user := models.User{}
	userId := utils.GetUserId(c)

	err := db.First(&user, userId).Error
	if err == gorm.ErrRecordNotFound {
		return utils.ReturnAlert(c, http.StatusNotFound, "not_found")
	} else if err != nil {
		return utils.ReturnAlert(c, http.StatusInternalServerError, "internal")
	}

	user.Email = utils.CreateSHA256(fmt.Sprint(userId)) + "+" + user.Email
	err = db.Save(&user).Error
	if err == gorm.ErrRecordNotFound {
		return utils.ReturnAlert(c, http.StatusNotFound, "not_found")
	} else if err != nil {
		return utils.ReturnAlert(c, http.StatusInternalServerError, "internal")
	}

	err = db.Delete(&models.User{}, userId).Error
	if err == gorm.ErrRecordNotFound {
		return utils.ReturnAlert(c, http.StatusNotFound, "not_found")
	} else if err != nil {
		return utils.ReturnAlert(c, http.StatusInternalServerError, "internal")
	}

	err = db.Delete(&models.Token{}, "user_id", userId).Error
	if err == gorm.ErrRecordNotFound {
		return utils.ReturnAlert(c, http.StatusNotFound, "not_found")
	} else if err != nil {
		return utils.ReturnAlert(c, http.StatusInternalServerError, "internal")
	}

	return c.JSON(http.StatusOK, map[string]any{"status": true, "data": []any{}})
}

func EditUserProfile(c echo.Context) error {
	var user models.User

	if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
		return utils.ReturnAlert(c, http.StatusBadRequest, "bad_request")
	}

	if utils.ValidatePassword(user.Password) {
		return utils.ReturnAlert(c, http.StatusBadRequest, "insecure_password")
	}
	if _, err := mail.ParseAddress(user.Email); err != nil {
		return utils.ReturnAlert(c, http.StatusBadRequest, "invalid_email")
	}

	user.Password = utils.CreateSHA256(user.Password)
	err := db.Where(utils.GetUserId(c)).Omit("role", "blocked", "verified", "user_id").Updates(&user).Error
	if err == gorm.ErrRecordNotFound {
		return utils.ReturnAlert(c, http.StatusNotFound, "not_found")
	} else if err != nil {
		return utils.ReturnAlert(c, http.StatusInternalServerError, "internal")
	}

	return c.JSON(http.StatusOK, map[string]any{"status": true, "data": []any{}})
}
