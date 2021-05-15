package controllers

import (
	"encoding/json"
	"github.com/KMSIGN/nsd-hack/MainServer/models"
	u "github.com/KMSIGN/nsd-hack/MainServer/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateMessage(c *gin.Context) {
	user := c.MustGet("user") //Grab the id of the user that send the request
	message := &models.Message{}

	err := json.NewDecoder(c.Request.Body).Decode(message)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"message":"Error while decoding request body"})
		return
	}

	message.UserId = int(user.(uint))
	resp := message.Create()
	c.JSON(http.StatusOK, resp)
}

func GetMessagesFor(c *gin.Context) {

	id := c.MustGet("user")
	data := models.GetMessages(int(id.(uint)))
	resp := u.Message(true, "success")
	resp["data"] = data
	c.JSON(http.StatusOK, resp)
}

