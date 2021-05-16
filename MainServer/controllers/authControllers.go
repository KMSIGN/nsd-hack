package controllers

import (
	"encoding/json"
	"github.com/KMSIGN/nsd-hack/MainServer/models"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
)

func CreateAccount(c *gin.Context) {

	account := &models.Account{}
	log.Println(c.Request.Body)
	err := json.NewDecoder(c.Request.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
		return
	}
	account.TelegToken = RandomString(30)

	resp := account.Create() //Create account
	c.JSON(http.StatusOK, resp)
}

func Authenticate(c *gin.Context) {
	account := &models.Account{}
	err := json.NewDecoder(c.Request.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
		return
	}

	resp := models.Login(account.Login, account.Password)
	c.JSON(http.StatusOK, resp)
}

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}