package routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
	UserController "summa-auth-api/controllers/v1"
	"summa-auth-api/controllers/v1"
)

func LoadV1(e *echo.Echo)  {
	//load Prisma
	v1.LoadPrisma()
	g := e.Group("/api/v1")
	g.GET("/", StatusCheckV1)
	g.GET("/test", UserController.StatusCheck)
	g.POST("/signup", UserController.CreateUser)
}

func StatusCheckV1(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string] interface{} {"status": true, "message": "Summa Auth API Version 1"})
}
