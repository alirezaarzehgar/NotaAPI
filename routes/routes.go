package routes

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	"github.com/Asrez/NotaAPI/config"
)

func todo(c echo.Context) error { return nil }

func Init() *echo.Echo {
	e := echo.New()
	e.POST("/register", todo)
	e.POST("/login", todo)
	e.POST("/token/create", todo)

	g := e.Group("", echojwt.WithConfig(echojwt.Config{SigningKey: config.JwtSecret()}))
	g.POST("/story/create", todo)
	return e
}
