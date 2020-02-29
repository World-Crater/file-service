package imageController

import (
	"os"
	"fmt"
	"bytes"

	"github.com/gin-gonic/gin"

	helperImage "file-service/helper"
	imgur "file-service/external"
	modelImage "file-service/model"
)

func Upload(c *gin.Context){
	getImageBuffer := func()(*bytes.Buffer){
		return helperImage.GetFormFile("upload", c)
	}
	errorHandler := func(err error, message string){
		fmt.Println(err)
		c.JSON(500, gin.H{
			"message": message,
		})
	}
	getImageSize := func()(int64){
		_, a, _ := c.Request.FormFile("upload")
		return a.Size
	}
	dynamodbFileServiceImageName := os.Getenv("DYNAMODB_FILE_SERVICE_IMAGES_NAME")
	var imageMaxSize int64 = 100

	width, _, err := helperImage.GetSize(getImageBuffer())
	if err != nil {
		c.Abort()
		errorHandler(err, "Image error")
		return
	}

	var processedImage *bytes.Buffer

	imageSize := getImageSize()
	if imageSize >= imageMaxSize {
		processedImage, err = helperImage.ImageResizeByBuffer(getImageBuffer(), width)
		if err != nil {
			c.Abort()
			errorHandler(err, "Image error")
			return
		}
	} else {
		processedImage = getImageBuffer()
	}

	imgurURL, err := imgur.Upload(processedImage)
	if err != nil {
		c.Abort()
		errorHandler(err, "Image error")
		return
	}
	_, err = modelImage.SaveImage(dynamodbFileServiceImageName, imgurURL)
	if err != nil {
		c.Abort()
		errorHandler(err, "Image error")
		return
	}
	c.JSON(200, gin.H{
		"imgurURL": imgurURL,
	})
}