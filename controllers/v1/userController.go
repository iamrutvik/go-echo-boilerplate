package v1

import (
	"github.com/dgrijalva/jwt-go"
	Struct "github.com/fatih/structs"
	"github.com/labstack/echo/v4"
	"net/http"
	Helpers "summa-auth-api/helpers"
	Models "summa-auth-api/models"
	"summa-auth-api/prisma-client"
	"time"
)

// statusCheck godoc
// @Summary ping example
// @Description do ping
// @Tags Status
// @Accept json
// @Produce json
// @Success 200 {string} string "pong"
// @Failure 400 {string} string "ok"
// @Failure 404 {string} string "ok"
// @Failure 500 {string} string "ok"
// @Router /test [get]
func StatusCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string] interface{} {"status": true, "message": "Hello world from v1"})
	//return echo.NewHTTPError(http.StatusInternalServerError,  "Hello world")
}
type (
	JSONResult struct {
		Status  bool              `json:"status" `
		Message string            `json:"message"`
		Data    interface{}       `json:"data"`
		Error   interface{}       `json:"error"`
	}
)
// CreateUser godoc
// @Summary Sign up
// @Description Create a new user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param User body models.User true "Add user"
// @Success 200 {object} JSONResult{data=models.User} "User Response, it will also return Token and hides Password and OTP"
// @Failure 400 {object} JSONResult{} "Validation error response with message"
// @Failure 500 {object} JSONResult{} "Internal Server error response with message"
// @Router /api/v1/signup [post]
func CreateUser(c echo.Context) error {
	u := new(Models.User)
	if err := c.Bind(u); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,  "Can not bind request body to type User")
	}
	if err := c.Validate(u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	var (
		isActive   = true
		isVerified = false
		otp        = Helpers.EncodeToString(6)
		ut         = prisma.UserType(u.UserType)
		ls         = prisma.LoginSource(u.LoginSource)
		bday       = "9999-12-31"
	)
	if u.Birthday != "" {
		bday = u.Birthday
	}
	user, err := Client.CreateUser(prisma.UserCreateInput{
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		Email:          u.Email,
		Mobile:         u.Mobile,
		Birthday:       &bday,
		Password:       Helpers.HashAndSalt(u.Password),
		ProfilePicture: u.ProfilePicture,
		UserType:       &ut,
		LoginSource:    &ls,
		IsVerified: 	&isVerified,
		IsActive:       &isActive,
		Otp:            &otp,
	}).Exec(Context)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	//load config
	Configuration, err := Helpers.LoadConfig()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set claims
	// This is the information which frontend can use
	// The backend can also decode the token and get admin etc.
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = user.FirstName + " " + user.LastName
	claims["email"] = user.Email
	claims["mobile"] = user.Mobile
	claims["userType"] = user.UserType
	claims["isVerified"] = user.IsVerified
	claims["isActive"] = user.IsActive
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(Configuration.Application.JWTExpireAt)).Unix()
	// Generate encoded token and send it as response.
	// The signing string should be secret (a generated UUID          works too)
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	userMap := Struct.Map(user)
	delete(userMap, "Password")
	delete(userMap, "Otp")
	userMap["Token"] = t
	return c.JSON(http.StatusOK, Helpers.Response(userMap))
}
