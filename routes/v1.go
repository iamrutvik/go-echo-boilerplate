package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	Config "summa-auth-api/config"
	"summa-auth-api/controllers/v1"
	UserController "summa-auth-api/controllers/v1"
)

func LoadV1(e *echo.Echo)  {
	//Load Prisma
	v1.LoadPrisma()

	g := e.Group("/api/v1")
	g.GET("/", StatusCheckV1)
	g.GET("/test", UserController.StatusCheck)
	g.POST("/signup", UserController.CreateUser)
	g.POST("/login", UserController.Login)

	userRoutes := e.Group("/api/v1/user")
	userRoutes.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &UserController.JWTCustomClaims{},
		SigningKey: []byte(Config.Settings.Application.JWTSecret),
	}))
	userRoutes.POST("/changepassword", UserController.ChangePassword)
}

func StatusCheckV1(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string] interface{} {"status": true, "message": "Summa Auth API Version 1"})
}
