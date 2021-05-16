package api

import (
	"fmt"
	"strconv"

	"github.com/KMSIGN/nsd-hack/server/app/file"
	"github.com/KMSIGN/nsd-hack/server/app/loader"

	"github.com/gin-gonic/gin"
)

func handleUpload(c *gin.Context) {
	hashes := c.PostForm("hashes")
	hash := c.PostForm("hash")
	lastPartSize, err := strconv.Atoi(c.PostForm("lastPartSize"))
	if err != nil {
		c.String(500, "Something go wrong: ", err)
		return
	}

	port, err := loader.SrvFileLoader(hash, hashes, lastPartSize)
	if err != nil {
		c.String(500, "Something go wrong: ", err)
		return
	}
	c.String(200, ":%d", port)
}

func handleDownload(c *gin.Context) {
	name := c.PostForm("hash")

	port, err := loader.StartUploading(name)
	if err != nil {
		c.String(500, "Something go wrong: ", err)
		return
	}

	c.String(200, ":%d", port)
}

func handleGetFilesCount(c *gin.Context) {
	c.String(200, fmt.Sprintf("%d", file.GetFilesLen()))
}
