package routes

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/Asrez/NotaAPI/api/handlers"
	"github.com/Asrez/NotaAPI/api/middlewares"
	"github.com/Asrez/NotaAPI/config"
)

func todo(c echo.Context) error { return nil }

func Init() *echo.Echo {
	e := echo.New()
	e.HTTPErrorHandler = handlers.DefaultLogHandler

	if config.Debug() {
		e.Static("/logs/", config.LogDirectory())
		e.GET("/logs/list", handlers.ShowLogs)
		e.GET("/logs/current", handlers.ShowCurrentLogs)
		e.Use(middlewares.LogRequests)
	}

	e.Use(echomiddleware.CORS())
	e.POST("/user/register", handlers.Register)
	e.POST("/user/login", handlers.Login)

	e.POST("/token/check", handlers.CheckToken)
	e.POST("/guest/create-token", handlers.CreateGuestToken)

	g := e.Group("", echojwt.WithConfig(echojwt.Config{SigningKey: config.JwtSecret()}), middlewares.CheckToken)
	g.Static("/", config.Assets())

	u := g.Group("", middlewares.UserOnly)
	u.GET("/user/story/count", handlers.GetStoryCount)
	u.GET("/user/profile", handlers.GetUserProfile)
	u.PUT("/user/profile", handlers.EditUserProfile)
	u.POST("/user/business-name-exists/:name", handlers.CheckBusinessNameExistance)
	u.DELETE("/user/delete-account", handlers.UserDeleteAccount)

	g.GET("/story/:code", handlers.GetStoryInfo)
	g.GET("/story/exists/:code", handlers.CheckStoryExistance)
	u.POST("/story/upload-asset", handlers.UploadAsset)
	u.POST("/story/name-exists/:name", handlers.CheckStoryNameExistance)
	u.POST("/story/create", handlers.CreateStory)
	u.POST("/story/change-status/:code", handlers.ChangeStoryStatus)
	u.GET("/story/stories", handlers.ListStories)
	u.PUT("/story/:code", handlers.EditStoryInfo)
	u.DELETE("/story/:code", handlers.DeleteStory)
	u.POST("/story/convert/:code", handlers.ConvertStory)

	g = g.Group("", middlewares.GuestOnly)
	g.GET("/guest/settings", handlers.GetGuestSettings)
	g.PUT("/guest/settings", handlers.EditGuestSettings)
	g.POST("/guest/save-story/:code", handlers.SaveStoryForGuest)
	g.GET("/guest/stories", handlers.ListGuestStories)
	g.GET("/guest/story/count", handlers.GetGuestStoryCount)
	g.GET("/guest/available-story-dates", handlers.ListStoryDates)
	g.GET("/guest/min-max-dates", handlers.GetMinAndMaxStoryDates)
	g.DELETE("/guest/delete-account", handlers.GuestDeleteAccount)

	return e
}
