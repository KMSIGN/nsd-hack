package api

import (
	"github.com/gin-gonic/gin"
	"nsd-hack/server/app/loader"
)

func handleUpload(c *gin.Context){
	hashes := c.PostForm("hashes")
	hash := c.PostForm("hash")

	port, err := loader.SrvFileLoader(hash, hashes)
	if err != nil {
		c.String(500, "Something go wrong")
	}
	c.String(200, ":%d", port)
}