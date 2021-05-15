package api

import (
	"github.com/KMSIGN/nsd-hack/server/app/loader"

	"github.com/gin-gonic/gin"
)

func handleUpload(c *gin.Context) {
	hashes := c.PostForm("hashes")
	hash := c.PostForm("hash")

	port, err := loader.SrvFileLoader(hash, hashes)
	if err != nil {
		c.String(500, "Something go wrong")
	}
	c.String(200, ":%d", port)
}
