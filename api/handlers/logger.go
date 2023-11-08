package handlers

import (
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"
	"time"

	"github.com/Asrez/NotaAPI/config"
	"github.com/labstack/echo/v4"
)

func ShowLogs(c echo.Context) error {
	list := ""
	filepath.Walk(config.LogDirectory(), func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		path = filepath.Base(path)
		list += fmt.Sprintf("<a href=\"%s\">%s</a>\n\n", path, path)
		return nil
	})
	return c.HTML(http.StatusOK, list)
}

func ShowCurrentLogs(c echo.Context) error {
	url := fmt.Sprintf("/logs/%s.log", time.Now().Format("2006-01-02"))
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func DefaultLogHandler(err error, c echo.Context) {
	he := err.(*echo.HTTPError)
	c.JSON(he.Code, map[string]any{
		"status": false,
		"alert":  he.Message,
		"data":   []any{},
	})
}
