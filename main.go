package main

import (
	imageController "file-service/controller"
	"log"

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
	r.POST("image/imgur", imageController.UploadToImgur)
	r.POST("image/s3", imageController.UploadToS3)
	r.Run()
}
