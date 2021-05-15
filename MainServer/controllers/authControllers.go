package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go-contacts/models"
	"log"
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
