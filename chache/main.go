	package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-contacts/app"
	"go-contacts/controllers"
	"net/http"
	"os"
)

func main() {

	router := gin.Default()

	auth := router.Group("/api/user")
	auth.POST("/new", controllers.CreateAccount)
	auth.POST("/login", controllers.Authenticate)

	mystat := router.Group("/api/files")
	mystat.Use(app.JwtAuthentication)
	mystat.POST("/new", controllers.CreateFile)
	mystat.GET("/list", controllers.GetFilesFor)

	messages := router.Group("/api/messages")
	messages.Use(app.JwtAuthentication)
	messages.POST("/new", controllers.CreateMessage)
	messages.GET("/", controllers.GetMessagesFor)


	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" //localhost
	}

	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}
