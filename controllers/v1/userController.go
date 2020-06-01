package v1

import (
	"github.com/dgrijalva/jwt-go"
	Struct "github.com/fatih/structs"
	"github.com/labstack/echo/v4"
	"net/http"
	Config "summa-auth-api/config"
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
	JWTCustomClaims struct {
		Name  string `json:"name"`
		Email string   `json:"email"`
		Mobile string   `json:"mobile"`
		UserType string `json:"userType"`
		IsVerified bool `json:"isVerified"`
		IsActive bool `json:"isActive"`
		jwt.StandardClaims
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
		isVerified = true //TODO verify the phone number
		otp        = Helpers.EncodeToString(6)
		ut         = prisma.UserType(u.UserType)
		ls         = prisma.LoginSource(u.LoginSource)
		bday       = "9999-12-31" //TODO dont want to give static bday
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

	// Create token
	// Set custom claims
	claims := &JWTCustomClaims{
		user.FirstName + " " + user.LastName,
		user.Email,
		user.Mobile,
		string(user.UserType),
		user.IsVerified,
		user.IsActive,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(Config.Settings.Application.JWTExpires)).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	// The signing string should be secret (a generated UUID          works too)
	t, err := token.SignedString([]byte(Config.Settings.Application.JWTSecret))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	userMap := Struct.Map(user)
	delete(userMap, "Password")
	delete(userMap, "Otp")
	userMap["Token"] = t
	return c.JSON(http.StatusOK, Helpers.Response(userMap))
}

// Login godoc
// @Summary Login
// @Description Login user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param User body models.LoginUser true "Add Login credentials"
// @Success 200 {object} JSONResult{data=models.User} "User Response, it will also return Token and hides Password and OTP"
// @Failure 400 {object} JSONResult{} "Validation error response with message"
// @Failure 500 {object} JSONResult{} "Internal Server error response with message"
// @Router /api/v1/login [post]
func Login(c echo.Context) error {
	u := new(Models.LoginUser)
	if err := c.Bind(u); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,  "Can not bind request body to type User")
	}
	if err := c.Validate(u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := Client.User(prisma.UserWhereUniqueInput{
			Email: &u.Email,
	}).Exec(Context)

	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	if Helpers.ComparePasswords(user.Password, u.Password) &&
		user.IsActive && user.IsVerified {
		// Create token
		// Set custom claims
		claims := &JWTCustomClaims{
			user.FirstName + " " + user.LastName,
			user.Email,
			user.Mobile,
			string(user.UserType),
			user.IsVerified,
			user.IsActive,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * time.Duration(Config.Settings.Application.JWTExpires)).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		// The signing string should be secret (a generated UUID          works too)
		t, err := token.SignedString([]byte(Config.Settings.Application.JWTSecret))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		userMap := Struct.Map(user)
		delete(userMap, "Password")
		delete(userMap, "Otp")
		userMap["Token"] = t
		return c.JSON(http.StatusOK, Helpers.Response(userMap))
	} else {
		return echo.NewHTTPError(http.StatusUnauthorized, "Password is wrong")
	}
}

// ChangePassword godoc
// @Summary Change Password
// @Description Login user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param User body models.LoginUser true "Add Login credentials"
// @Success 200 {object} JSONResult{data=models.User} "User Response, it will also return Token and hides Password and OTP"
// @Failure 400 {object} JSONResult{} "Validation error response with message"
// @Failure 500 {object} JSONResult{} "Internal Server error response with message"
// @Router /api/v1/login [post]
func ChangePassword(c echo.Context) error {
	u := new(Models.ChangePassword)
	if err := c.Bind(u); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,  "Can not bind request body to type User")
	}
	if err := c.Validate(u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	authUser := c.Get("user").(*jwt.Token)
	claims := authUser.Claims.(*JWTCustomClaims)
	email := claims.Email
	hashedPassword := Helpers.HashAndSalt(u.Password)

	user, err := Client.UpdateUser(prisma.UserUpdateParams{
		Where: prisma.UserWhereUniqueInput{
			Email: &email,
		},
		Data: prisma.UserUpdateInput{
			Password: &hashedPassword,
		},
	}).Exec(Context)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	userMap := Struct.Map(user)
	delete(userMap, "Password")
	delete(userMap, "Otp")
	return c.JSON(http.StatusOK, Helpers.Response(userMap))
}


