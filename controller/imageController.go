package imageController

import (
	"bytes"
	"fmt"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"

	imgur "file-service/external"
	helperImage "file-service/helper"
	modelImage "file-service/model"
)

func UploadToImgur(c *gin.Context) {
	getImageBuffer := func() *bytes.Buffer {
		return helperImage.GetFormFile("upload", c)
	}
	errorHandler := func(err error, message string) {
		fmt.Println(err)
		c.JSON(500, gin.H{
			"message": message,
		})
	}
	getImageSize := func() int64 {
		_, a, _ := c.Request.FormFile("upload")
		return a.Size
	}
	dynamodbFileServiceImageName := viper.GetString("DYNAMODB_FILE_SERVICE_IMAGES_NAME")
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
		"url": imgurURL,
	})
}

func UploadToS3(c *gin.Context) {
	getImageBuffer := func() *bytes.Buffer {
		return helperImage.GetFormFile("upload", c)
	}
	errorHandler := func(err error, message string) {
		fmt.Println(err)
		c.JSON(500, gin.H{
			"message": message,
		})
	}
	getImageSize := func() int64 {
		_, a, _ := c.Request.FormFile("upload")
		return a.Size
	}
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
	uuid4 := uuid.Must(uuid.NewV4())
	keyName := uuid4.String() + ".jpg"
	_, err = modelImage.SaveImageToS3(processedImage, keyName)
	if err != nil {
		c.Abort()
		errorHandler(err, "Image upload fail")
		return
	}
	c.JSON(200, gin.H{
		"url": "https://s3-" + viper.GetString("AWS_REGION") + ".amazonaws.com/" + viper.GetString("IMAGE_BUCKET_NAME") + "/" + keyName,
	})
}
