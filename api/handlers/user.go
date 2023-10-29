package handlers

import (
	"encoding/json"
	"net/http"
	"net/mail"

	"github.com/labstack/echo/v4"

	"github.com/Asrez/NotaAPI/models"
	"github.com/Asrez/NotaAPI/utils"
)

func Register(c echo.Context) error {
	user := models.User{Role: models.USERS_ROLE_USER}
	if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
		return echo.ErrBadRequest
	}
	if user.Email == "" || user.Password == "" {
		return echo.ErrBadRequest
	}

	if utils.ValidatePassword(user.Password) {
		return utils.ReturnAlert(c, http.StatusBadRequest, "insecure_password")
	}
	if _, err := mail.ParseAddress(user.Email); err != nil {
		return utils.ReturnAlert(c, http.StatusBadRequest, "invalid_email")
	}

	user.Password = utils.HashPassword(user.Password)

	r := db.Create(&user)
	switch {
	case r.Error != nil:
		return utils.ReturnAlert(c, http.StatusInternalServerError, "internal")
	case r.RowsAffected == 0:
		return utils.ReturnAlert(c, http.StatusConflict, "user_conflict")
	}

	token := utils.CreateUserToken(user.ID, user.Email, user.Username)
	if err := db.Create(&models.Token{JwtToken: token}).Error; err != nil {
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
		return echo.ErrBadRequest
	}
	fillteredUser := models.User{Email: user.Email, Password: utils.HashPassword(user.Password)}
	db.Where(fillteredUser).First(&user).Count(&loggedin)

	if loggedin == 0 {
		return utils.ReturnAlert(c, http.StatusUnauthorized, "login_unauthorized")
	}

	token := utils.CreateUserToken(user.ID, user.Email, user.Username)
	if err := db.Create(&models.Token{JwtToken: token}).Error; err != nil {
		return utils.ReturnAlert(c, http.StatusInternalServerError, "internal")
	}

	return c.JSON(http.StatusOK, map[string]any{
		"status": true,
		"data":   map[string]any{"token": token},
	})
}
