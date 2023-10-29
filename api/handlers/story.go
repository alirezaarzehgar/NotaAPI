package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/Asrez/NotaAPI/config"
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
