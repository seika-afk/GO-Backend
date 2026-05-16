package controllers

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetUsers() gin.HandlerFunc {

	return func(c *gin.Context) {

	}
}

func GetUser() gin.HandlerFunc {

	return func(c *gin.Context) {

	}

}

func Signup() gin.HandlerFunc {

	return func(c *gin.Context) {

	}
}

func Login() gin.HandlerFunc {

	return func(c *gin.Context) {

	}

}

func HashPassword(psd string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(psd), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(bytes)
}

// userPassword and providedPassword
func VerifyPassword(uP string, pP string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(uP), []byte(pP))
	if err != nil {
		return false, err.Error()
	}
	return true, ""
}
