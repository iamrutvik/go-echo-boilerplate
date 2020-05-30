package routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
	UserController "summa-auth-api/controllers/v2"
)

func LoadV2(e *echo.Echo)  {
	g := e.Group("/api/v2")
	g.GET("/", StatusCheckV2)
	g.GET("/test", UserController.StatusCheck)
}

func StatusCheckV2(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string] interface{} {"status": true, "message": "Summa Auth API Version 2"})
}
