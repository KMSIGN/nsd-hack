package api

import (
	"fmt"
	"github.com/KMSIGN/nsd-hack/server/app/file"
	"github.com/KMSIGN/nsd-hack/server/app/loader"

	"github.com/gin-gonic/gin"
)

func handleUpload(c *gin.Context) {
	hashes := c.PostForm("hashes")
	hash := c.PostForm("hash")

	port, err := loader.SrvFileLoader(hash, hashes)
	if err != nil {
		c.String(500, "Something go wrong: ", err)
		return
	}
	c.String(200, ":%d", port)
}

func handleDownload(c *gin.Context){
	addr := c.PostForm("addr")
	name := c.PostForm("hash")

	err := loader.StartUploading(addr, name)
	if err != nil {
		c.String(500, "Something go wrong: ", err)
		return
	}

	c.String(200, "ok")
}

func handleGetFilesCount(c *gin.Context){
	c.String(200, fmt.Sprintf("%d", file.GetFilesLen()))
}
