package main

import (
	imageController "file-service/controller"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	// init
	viper.AutomaticEnv()

	// route
	r := gin.Default()
	r.POST("image/imgur", imageController.UploadToImgur)
	r.POST("image/s3", imageController.UploadToS3)
	r.Run()
}
