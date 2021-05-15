package scheduler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func RederictMidlerare(addr string) gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.FullPath() == "/api/file/new" {
			var temp map[string]interface{}

			err := json.NewDecoder(c.Request.Body).Decode(&temp)
			if err != nil {
				c.JSON(http.StatusBadRequest, map[string]string{"message":"Error while decoding request body"})
				return
			}

			temp["cache"] = os.Getenv("fileserver")
			//temp["timeout"] =
		}

		c.Redirect(http.StatusMovedPermanently, addr + c.FullPath())
	}
}
