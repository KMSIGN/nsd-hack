package main

import (
	"flag"
	"github.com/KMSIGN/nsd-hack/CacheServer/scheduler"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func init() {
	e := godotenv.Load()
	if e != nil {
		log.Fatal(e)
	}
}

func main(){

	mainS := flag.String("addr", "", "Main server addr")
	if *mainS == "" {
		log.Fatal("Empty main server address")
	}

	router := gin.Default()
	router.Group("/").Use(scheduler.RederictMidlerare(*mainS))
}
