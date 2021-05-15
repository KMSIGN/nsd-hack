package controllers

import (
"encoding/json"
"github.com/gin-gonic/gin"
"go-contacts/models"
u "go-contacts/utils"
"net/http"
)

func CreateMessage(c *gin.Context) {
	user := c.GetInt("user") //Grab the id of the user that send the request
	message := &models.Message{}

	err := json.NewDecoder(c.Request.Body).Decode(message)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"message":"Error while decoding request body"})
		return
	}

	message.UserId = user
	resp := message.Create()
	c.JSON(http.StatusOK, resp)
}

func GetMessagesFor(c *gin.Context) {

	id := c.GetInt("user")
	data := models.GetMessages(id)
	resp := u.Message(true, "success")
	resp["data"] = data
	c.JSON(http.StatusOK, resp)
}

