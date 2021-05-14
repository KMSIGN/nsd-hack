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
	router.GET("/upload", handleUpload)
}

func ListenAndServe() error {
	return http.ListenAndServe(":8080", router)
}