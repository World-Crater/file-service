package main

import (
	"log"
	"file-service/controller"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// init
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// route
	r := gin.Default()
	r.POST("image/upload", imageController.Upload)
	r.Run()
}