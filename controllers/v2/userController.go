package v2

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func StatusCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string] interface{} {"status": true, "message": "Hello world from v2"})
	//return echo.NewHTTPError(http.StatusInternalServerError,  "Hello world")
}
