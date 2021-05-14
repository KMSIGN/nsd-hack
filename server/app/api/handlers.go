package api

import (
	"github.com/gin-gonic/gin"
	"nsd-hack/server/app/loader"
)

func handleUpload(c *gin.Context){
	hashes := c.PostForm("hashes")
	name := nameHashing(hashes)

	port, err := loader.SrvFileLoader(name, hashes)
	if err != nil {
		c.String(500, "Something go wrong")
	}
	c.String(200, ":%d", port)
}