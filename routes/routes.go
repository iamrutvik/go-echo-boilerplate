package routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"summa-auth-api/config"
)

func Load(configuration config.Configurations, e *echo.Echo){
	e.GET("/", func (c echo.Context) error {
		return c.JSON(http.StatusOK, map[string] interface{} {"status": true, "message": configuration.Application.Name + " running"})
	})

	//add index function for each API version in below slice
	for _, fn := range []func(e *echo.Echo){LoadV1, LoadV2} {
		fn(e)
	}

}
