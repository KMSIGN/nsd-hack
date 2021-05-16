package controllers

import (
	"github.com/KMSIGN/nsd-hack/MainServer/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func StatisticToken(c *gin.Context){
	user := c.MustGet("user") //Grab the id of the user that send the request

	ur := models.GetUser(user.(uint))
	if ur == nil {
		c.JSON(http.StatusBadRequest, map[string]string{"message":"User not found"})
		return
	}

	data := map[string]string{
		"status": "ok",
		"token": ur.TelegToken,
	}

	c.JSON(http.StatusOK, data)
}

func GetUserStatistic(c *gin.Context){
	user := c.GetInt("user")

	data := map[string]string{
		"FileSumSize":strconv.Itoa(models.GetFilesSizeSum(user)),
		"TotalFileCount": strconv.Itoa(models.GetFilesCount(user)),
	}
	c.JSON(http.StatusOK, data)
}
