package helpers

import (
	"crypto/rand"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"strconv"
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

func ComparePasswords(hashedPwd string, plainPwd string) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	bytePlainPwd := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePlainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func Response (res map[string]interface{}) JSONResult {
	return JSONResult{
		Status:  true,
		Message: "Success",
		Data:    res,
		Error:   nil,
	}
}