package helpers

import (
	"crypto/rand"
	"errors"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"strconv"
	"summa-auth-api/config"
)

type (
	JSONResult struct {
		Status  bool                   `json:"status" `
		Message string                 `json:"message"`
		Data    map[string]interface{} `json:"data"`
		Error   interface{}            `json:"error"`
	}
)

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
func EncodeToString(max int) int32 {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	otp, _ := strconv.Atoi(string(b))
	return int32(otp)
}

func HashAndSalt(pwd string) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

//TODO: Think about loading config once globally for app
func LoadConfig() (config.Configurations, error) {
	// Start viper implementation - reading configurations from configurations.yaml
	// Set the file name of the configurations file
	viper.SetConfigName("config")

	// Set the path to look for the configurations file
	viper.AddConfigPath(".")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType("yaml")
	var configurations config.Configurations

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return configurations, errors.New("configuration file not found at ./config/config.yml")
		} else {
			return configurations, errors.New("error reading configurations file")
		}
	}

	//setting defaults in case of undefined
	viper.SetDefault("server.port", "8000")
	viper.SetDefault("server.certfile", "./cert.pem")
	viper.SetDefault("server.keyfile", "./key.pem")

	//Decoding viper configurations to struct configuration from ./configurations/configurations.go
	err := viper.Unmarshal(&configurations)
	if err != nil {
		return configurations, errors.New("unable to decode into struct")
	}

	return configurations, nil
	//End Viper implementation
}

func Response (res map[string]interface{}) JSONResult {
	return JSONResult{
		Status:  true,
		Message: "Success",
		Data:    res,
		Error:   nil,
	}
}