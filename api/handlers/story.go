package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

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
	utils.DebugLog("start saving an asset")

	file, err := c.FormFile("asset")
	if err != nil {
		return utils.ReturnAlert(c, http.StatusBadRequest, "bad_request")
	}
	utils.DebugLog("recieve file:", file.Filename)

	src, err := file.Open()
	if err != nil {
		log.Println("open file : ", err)
		return utils.ReturnAlert(c, http.StatusInternalServerError, "internal")
	}
	defer src.Close()
	utils.DebugLog("open sent file:", file.Filename)

	if !utils.IsValidPath(file.Filename, isImage) {
		return utils.ReturnAlert(c, http.StatusBadRequest, "bad_file")
	}

	dirpath := utils.GetUserDir(utils.GetUserId(c))
	assetsDir := config.Assets() + "/" + dirpath
	if _, err := os.Stat(assetsDir); err != nil {
		if err := os.Mkdir(assetsDir, os.ModePerm); err != nil {
			log.Println("mkdir dirpath: ", err)
			return utils.ReturnAlert(c, http.StatusInternalServerError, "internal")
		}
		utils.DebugLog("create dir:", assetsDir)
	}

	filepath := fmt.Sprintf("%s/%s", dirpath, utils.GetUniqueName(file.Filename))
	dst, err := os.Create(config.Assets() + "/" + filepath)
	if err != nil {
		utils.DebugLog("create ", filepath, ": ", err)
		return utils.ReturnAlert(c, http.StatusInternalServerError, "internal")
	}
	defer dst.Close()
	utils.DebugLog("create asset on path:", filepath)

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	utils.DebugLog("copy transfered file to assets directory")

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
		return utils.ReturnAlert(c, http.StatusInternalServerError, "internal")
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
		return utils.ReturnAlert(c, http.StatusNotFound, "not_found")
	}
	return c.JSON(http.StatusOK, map[string]any{
		"status": true,
		"data":   map[string]any{"is_public": data.IsPublic},
	})
}

func CheckStoryExistance(c echo.Context) error {
	var count int64

	err := db.Model(&models.Story{}).
		Where(models.Story{Code: c.Param("code")}).
		Count(&count).Error
	if err != nil {
		return utils.ReturnAlert(c, http.StatusInternalServerError, "internal")
	}

	if count == 0 {
		return utils.ReturnAlert(c, http.StatusNotFound, "not_found")
	}

	return c.JSON(http.StatusOK, map[string]any{"status": true, "data": []any{}})
}

func ListStories(c echo.Context) error {
	var stories []models.Story
	dateCond := db.Where("1 = 1")
	defaultCond := map[string]any{"user_id": utils.GetUserId(c)}

	if c.QueryParam("story_type") == models.STORY_TYPE_EXPLORE {
		defaultCond["type"] = models.STORY_TYPE_EXPLORE
	} else if c.QueryParam("start_date") != "" && c.QueryParam("end_date") != "" {
		startDate, err := time.Parse(time.DateOnly, c.QueryParam("start_date"))
		if err != nil {
			return utils.ReturnAlert(c, http.StatusBadRequest, "bad_request")
		}
		endDate, err := time.Parse(time.DateOnly, c.QueryParam("end_date"))
		if err != nil {
			return utils.ReturnAlert(c, http.StatusBadRequest, "bad_request")
		}

		dateCond = db.Where("`to` >= ? AND `to` <= ?", startDate, endDate)
	}

	if v, err := strconv.ParseBool(c.QueryParam("is_public")); err == nil {
		defaultCond["is_public"] = v
	}

	r := db.Where(dateCond).Where(defaultCond).Find(&stories)
	if r.Error == gorm.ErrRecordNotFound {
		return utils.ReturnAlert(c, http.StatusNotFound, "not_found")
	}
	if r.Error != nil {
		return utils.ReturnAlert(c, http.StatusInternalServerError, "internal")
	}

	return c.JSON(http.StatusOK, map[string]any{
		"status": true,
		"data":   map[string]any{"stories": stories},
	})
}

func GetStoryInfo(c echo.Context) error {
	var story models.Story
	isUser := false
	userId := utils.GetUserId(c)
	code := c.Param("code")

	if userId > 0 {
		isUser = true
	}

	r := db.Where(models.Story{Code: code}).
		Or(map[string]any{"user_id": userId, "code": code}).First(&story)
	if r.Error == gorm.ErrRecordNotFound {
		return utils.ReturnAlert(c, http.StatusNotFound, "not_found")
	}
	if r.Error != nil {
		return utils.ReturnAlert(c, http.StatusInternalServerError, "internal")
	}

	if isUser && story.UserID == userId {
		return c.JSON(http.StatusOK, map[string]any{
			"status": true,
			"data":   map[string]any{"story": story},
		})
	}

	return c.JSON(http.StatusOK, map[string]any{"status": true, "data": story})
}

func EditStoryInfo(c echo.Context) error {
	var story models.Story

	if err := json.NewDecoder(c.Request().Body).Decode(&story); err != nil {
		return utils.ReturnAlert(c, http.StatusBadRequest, "bad_request")
	}

	err := db.Where(models.Story{UserID: utils.GetUserId(c), Code: c.Param("code")}).
		Omit("code", "is_public", "user_id").
		Updates(&story).Error
	if err == gorm.ErrRecordNotFound {
		return utils.ReturnAlert(c, http.StatusNotFound, "not_found")
	} else if err != nil {
		return utils.ReturnAlert(c, http.StatusInternalServerError, "internal")
	}

	return c.JSON(http.StatusOK, map[string]any{"status": true, "data": []any{}})
}

func DeleteStory(c echo.Context) error {
	story := models.Story{}
	err := db.First(&story, "code", c.Param("code")).Error
	if err == gorm.ErrRecordNotFound {
		return utils.ReturnAlert(c, http.StatusNotFound, "not_found")
	} else if err != nil {
		return utils.ReturnAlert(c, http.StatusInternalServerError, "internal")
	}

	err = db.Model(&story).Association("Tokens").Clear()
	if err != nil {
		return utils.ReturnAlert(c, http.StatusInternalServerError, "internal")
	}

	err = db.Delete(&story).Error
	if err != nil {
		return utils.ReturnAlert(c, http.StatusInternalServerError, "internal")
	}

	return c.JSON(http.StatusOK, map[string]any{"status": true, "data": []any{}})
}

func ConvertStory(c echo.Context) error {
	story := models.Story{}

	err := db.First(&story, "code = ? AND is_public = ?", c.Param("code"), true).Error
	if err == gorm.ErrRecordNotFound {
		return utils.ReturnAlert(c, http.StatusNotFound, "not_found")
	} else if err != nil {
		return utils.ReturnAlert(c, http.StatusInternalServerError, "internal")
	}
	if story.Type == models.STORY_TYPE_NORMAL {
		story.From = nil
		story.To = nil
		story.Type = models.STORY_TYPE_EXPLORE
	} else {
		var data struct {
			From time.Time `json:"from"`
			To   time.Time `json:"to"`
		}

		if err := json.NewDecoder(c.Request().Body).Decode(&data); err != nil {
			return utils.ReturnAlert(c, http.StatusBadRequest, "bad_request")
		}

		story.From = &data.From
		story.To = &data.To
		story.Type = models.STORY_TYPE_NORMAL
	}

	if field := validations.GetWrongStoryField(story); field != "" {
		return utils.ReturnAlert(c, http.StatusBadRequest, "story_wrong", field)
	}

	err = db.Save(&story).Error
	if err != nil {
		return utils.ReturnAlert(c, http.StatusInternalServerError, "internal")
	}

	return c.JSON(http.StatusOK, map[string]any{"status": true, "data": []any{}})
}
