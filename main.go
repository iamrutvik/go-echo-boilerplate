package main

import (
	"bytes"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"github.com/swaggo/echo-swagger" // echo-swagger middleware
	"io"
	"net/http"
	"strconv"
	Config "summa-auth-api/config"
	_ "summa-auth-api/docs"
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
	// Start viper implementation - reading configurations from configurations.yaml
	// Set the file name of the configurations file
	viper.SetConfigName("config")

	// Set the path to look for the configurations file
	viper.AddConfigPath(".")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType("yaml")


	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			e.Logger.Fatalf("Config file not found at root ./config.yml")
		} else {
			e.Logger.Fatalf("Error reading config file at root ./config.yml")
		}
	}

	//setting defaults in case of undefined
	viper.SetDefault("server.port", "8000")
	viper.SetDefault("server.certfile", "./cert.pem")
	viper.SetDefault("server.keyfile", "./key.pem")

	//Decoding viper configurations to struct configuration from ./configurations/configurations.go
	err := viper.Unmarshal(&Config.Settings)
	if err != nil {
		e.Logger.Fatalf("Unable to decode into struct")
	}
	//End Viper implementation

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
	Router.Load(e)
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
	startServer(e)
}

func startServer(e *echo.Echo) {
	if Config.Settings.Server.TLS {
		e.Logger.Infof("Server running on https://localhost:%v 🐹", Config.Settings.Server.Port)
		e.Logger.Fatal(
			e.StartTLS(":" + strconv.Itoa(Config.Settings.Server.Port),
				Config.Settings.Server.CertFile,
				Config.Settings.Server.KeyFile))
	} else {
		e.Logger.Infof("Server running on http://localhost:%v 🐹", Config.Settings.Server.Port)
		e.Logger.Fatal(e.Start(":" + strconv.Itoa(Config.Settings.Server.Port)))
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


