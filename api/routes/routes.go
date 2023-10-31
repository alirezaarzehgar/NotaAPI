package routes

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	"github.com/Asrez/NotaAPI/api/handlers"
	"github.com/Asrez/NotaAPI/api/middlewares"
	"github.com/Asrez/NotaAPI/config"
)

func todo(c echo.Context) error { return nil }

func Init() *echo.Echo {
	e := echo.New()
	e.POST("/user/register", handlers.Register)
	e.POST("/user/login", handlers.Login)

	e.POST("/token/check", handlers.CheckToken)
	e.POST("/guest/create-token", handlers.CreateGuestToken)

	g := e.Group("", echojwt.WithConfig(echojwt.Config{SigningKey: config.JwtSecret()}), middlewares.CheckToken)
	g.Static("/", config.Assets())
	g.GET("/user/story/count", handlers.GetStoryCount)
	g.DELETE("/user/delete-account", todo)

	g.GET("/story/:code", handlers.GetStoryInfo)
	g.GET("/story/exists/:code", handlers.CheckStoryExistance)
	u := g.Group("", middlewares.UserOnly)
	u.POST("/story/upload-asset", handlers.UploadAsset, middlewares.UserOnly)
	u.POST("/story/create", handlers.CreateStory)
	u.POST("/story/change-status/:code", handlers.ChangeStoryStatus)
	u.GET("/story/stories", handlers.ListStories)
	u.PUT("/story/:code", handlers.EditStoryInfo)
	u.DELETE("/story/:code", todo)
	u.POST("/story/convert", todo)

	g = g.Group("", middlewares.GuestOnly)
	g.GET("/guest/settings", handlers.GetGuestSettings)
	g.PUT("/guest/settings", handlers.EditGuestSettings)
	g.POST("/guest/save-story/:code", handlers.SaveStoryForGuest)
	g.GET("/guest/stories", handlers.ListGuestStories)
	g.DELETE("/guest/delete-account", handlers.GuestDeleteAccount)

	return e
}
