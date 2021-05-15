package app

import (
	"github.com/KMSIGN/nsd-hack/MainServer/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func TelegramToken(c *gin.Context) {
	tokenHeader := c.GetHeader("Token") //Grab the token from the header

	if tokenHeader == "" {
		c.Status(http.StatusForbidden)
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusForbidden, map[string]string{"message":"Missing telegram token"})
		c.Abort()
		return
	}

	user := models.GetUserByToken(tokenHeader)
	if user == nil {
		c.Status(http.StatusForbidden)
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusForbidden, map[string]string{"message":"Invalid telegram token"})
		c.Abort()
		return
	}

	c.Set("user", user.ID)
	c.Next()
}
