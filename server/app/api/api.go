package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	router *gin.Engine
)

func Configure(){
	router = gin.Default()
	router.POST("/upload", handleUpload)
	router.POST("/download", handleDownload)
	router.GET("/filesCount", handleGetFilesCount)
}

func ListenAndServe() error {
	return http.ListenAndServe(":8080", router)
}