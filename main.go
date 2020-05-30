package main

import (
	"bytes"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/swaggo/echo-swagger" // echo-swagger middleware
	"io"
	"net/http"
	"strconv"
	"summa-auth-api/config"
	_ "summa-auth-api/docs"
	Helpers "summa-auth-api/helpers"
	Router "summa-auth-api/routes"
)
// @title Summa Auth API
// @version 1.0
// @description Auth Service.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000
// @BasePath /api/v1

type CustomValidator struct {
	validator *validator.Validate
}

func main() {
	// Echo instance
	e := echo.New()

	//hide start up banner
	e.HideBanner = true

	//customErrorHandler
	e.HTTPErrorHandler = customErrorHandler

	//Logger config works only with e.Logger
	if l, ok := e.Logger.(*log.Logger); ok {
		l.SetHeader("${prefix} - ${level} => ")
		l.SetLevel(log.DEBUG)
		l.SetPrefix("Summa Auth API")
		l.EnableColor()
	}

	//Load config
	Configuration, err := Helpers.LoadConfig()
	if err != nil {
		e.Logger.Fatalf("Error loading Config - %s", err)
	}

	//can be done via `log` as well. log.SetLevel(log.DEBUG)
	//log.Infof("Server running on http://localhost%s 🐹", ":4000")

	// Middleware
	//e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
	//	Format: "method=${method}, uri=${uri}${path}, header=${header}, query=${query}, form=${form}, status=${status}\n",
	//}))
	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())
	e.Use(middleware.Secure())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH,    echo.POST, echo.DELETE},
	}))
	//e.Use(AfterResponse)

	//custom validator
	e.Validator = &CustomValidator{validator: validator.New()}

	// loads Routes
	Router.Load(Configuration, e)
	//Loading Swagger Route
	e.GET("/docs/*", echoSwagger.WrapHandler)

	//Start server
	/*
		Run the following command to generate cert.pem and key.pem files:
		go run $GOROOT/src/crypto/tls/generate_cert.go --host localhost
		Once the cert.pem is generated, double click and install it in macOS.
		It should be installed under "System" certificates in the left pane with name "Acme Co"
		Click on it, Go to Trust and Select Trust always
	 */
	startServer(Configuration, e)
}

func startServer(configuration config.Configurations, e *echo.Echo) {
	if configuration.Server.TLS {
		e.Logger.Infof("Server running on https://localhost:%v 🐹", configuration.Server.Port)
		e.Logger.Fatal(
			e.StartTLS(":" + strconv.Itoa(configuration.Server.Port),
				configuration.Server.CertFile,
				configuration.Server.KeyFile))
	} else {
		e.Logger.Infof("Server running on http://localhost:%v 🐹", configuration.Server.Port)
		e.Logger.Fatal(e.Start(":" + strconv.Itoa(configuration.Server.Port)))
	}
}

func customErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	var message interface{}
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = he.Message
	}
	NewMessage := map[string] interface{} {
		"status": false,
		"message": "Something went wrong",
		"error" : message,
		"data" : nil,
	}
	c.Logger().Error(err)
	_ = c.JSON(code, NewMessage)
}

//TODO implement a function to modify the response via middleware
func AfterResponse(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Logger().Infof("After called")
		c.Response().After(func() {
			res := c.Response()
			rw := res.Writer
			buf := new(bytes.Buffer)
			buf.WriteString(buf.String() + "{'done':'dome'}")
			io.MultiWriter(rw, buf)
			res.Writer = rw
			body := buf.String()
			c.Logger().Info(body)
			//err = next(c)
		})
		return next(c)
	}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}


