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
	g.GET("/user/story/count", todo)
	g.DELETE("/user/delete-account", todo)

	g.GET("/story/pub/:code", todo)
	g.POST("/story/create", todo)
	g.GET("/story/stories", todo)
	g.DELETE("/story/:code", todo)
	g.PUT("/story/:code", todo)
	g.POST("/story/publish/:code", todo)
	g.POST("/story/convert", todo)

	g.GET("/guest/settings", todo)
	g.PUT("/guest/settings", todo)
	g.DELETE("/guest/delete-account", todo)
	g.POST("/guest/save-story", todo)
	g.GET("/guest/stories", todo)

	return e
}
