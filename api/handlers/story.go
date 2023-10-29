package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/Asrez/NotaAPI/api/validations"
	"github.com/Asrez/NotaAPI/config"
	"github.com/Asrez/NotaAPI/models"
	"github.com/Asrez/NotaAPI/utils"
)

func UploadAsset(c echo.Context) error {
	var isImage bool
	if c.QueryParam("is_image") != "" {
		var err error
		isImage, err = strconv.ParseBool(c.QueryParam("is_image"))
		if err != nil {
			return utils.ReturnAlert(c, http.StatusBadRequest, "bad_request")
		}
	}

	file, err := c.FormFile("asset")
	if err != nil {
		return utils.ReturnAlert(c, http.StatusBadRequest, "bad_request")
	}
	src, err := file.Open()
	if err != nil {
		log.Println("open file : ", err)
		return utils.ReturnAlert(c, http.StatusInternalServerError, "internal")
	}
	defer src.Close()

	if !utils.IsValidPath(file.Filename, isImage) {
		return utils.ReturnAlert(c, http.StatusBadRequest, "bad_file")
	}

	dirpath := utils.GetUserDir(utils.GetUserId(c))
	if _, err := os.Stat(dirpath); err != nil {
		if err := os.Mkdir(dirpath, os.ModePerm); err != nil {
			log.Println("mkdir dirpath: ", err)
			return utils.ReturnAlert(c, http.StatusInternalServerError, "internal")
		}
	}

	filepath := fmt.Sprintf("%s/%s", dirpath, utils.GetUniqueName(file.Filename))
	dst, err := os.Create(config.Assets() + "/" + filepath)
	if err != nil {
		log.Println("create ", filepath, ": ", err)
		return utils.ReturnAlert(c, http.StatusInternalServerError, "internal")
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]any{
		"status": true,
		"data":   map[string]any{"path": filepath},
	})
}

func CreateStory(c echo.Context) error {
	var story models.Story

	if err := json.NewDecoder(c.Request().Body).Decode(&story); err != nil {
		return utils.ReturnAlert(c, http.StatusBadRequest, "bad_request", ":", err.Error())
	}

	if story.Type == models.STORY_TYPE_EXPLORE {
		var storyExploreCount int64

		err := db.Model(&models.Story{}).
			Where(models.Story{UserID: utils.GetUserId(c), Type: models.STORY_TYPE_EXPLORE}).
			Count(&storyExploreCount).Error
		if err != nil {
			return utils.ReturnAlert(c, http.StatusInternalServerError, "internal")
		}

		if storyExploreCount > 0 {
			return utils.ReturnAlert(c, http.StatusBadRequest, "dup_estory")
		}
	}

	if field := validations.GetWrongStoryField(story); field != "" {
		return utils.ReturnAlert(c, http.StatusBadRequest, "story_wrong", field)
	}

	story.UserID = utils.GetUserId(c)
	story.Code = utils.CreateRandomString(fmt.Sprint(story), 5)
	if err := db.Create(&story).Error; err != nil {
		return utils.ReturnAlert(c, http.StatusInternalServerError, "internal", ": create story: ", err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{
		"status": true,
		"data":   map[string]any{"code": story.Code},
	})
}

func ChangeStoryStatus(c echo.Context) error {
	var data struct {
		IsPublic bool `json:"is_public"`
	}
	if err := json.NewDecoder(c.Request().Body).Decode(&data); err != nil {
		return utils.ReturnAlert(c, http.StatusBadRequest, "bad_request", ":", err.Error())
	}

	r := db.Model(&models.Story{}).Where(models.Story{Code: c.Param("code")}).Update("is_public", data.IsPublic)
	if r.RowsAffected == 0 {
		return utils.ReturnAlert(c, http.StatusNotFound, "not_found", r.Error.Error())
	}
	return c.JSON(http.StatusOK, map[string]any{
		"status": true,
		"data":   map[string]any{"is_public": data.IsPublic},
	})
}
