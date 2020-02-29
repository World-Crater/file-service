package image

import (
	"bytes"
	"io"
	"log"
	"math"
	"image"
	_ "image/jpeg"
	"image/jpeg"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
)

func GetFormFile(formFilename string, c *gin.Context) (*bytes.Buffer){
	file, _, err := c.Request.FormFile("upload")
	imageBuffer := bytes.NewBuffer(nil)
	_, err = io.Copy(imageBuffer, file)
	if err != nil {
		log.Fatal(err)
	}
	return imageBuffer
}

func GetSize(imageBuffer *bytes.Buffer) (uint, uint, error) {
    m, _, err := image.Decode(imageBuffer)
    if err != nil {
		return 0, 0, err
    }
    g := m.Bounds()
    height := g.Dy()
    width := g.Dx()
	return uint(width), uint(height), nil
}

func ImageResizeByBuffer(file *bytes.Buffer, width uint) (*bytes.Buffer, error){
	base := file.Len() / 1000
	img, err := imaging.Decode(file)
	if err != nil {
		return nil, err
	}
	img = imaging.Resize(img, int(width / uint(math.Sqrt(float64(base))) * uint(math.Sqrt(10))), 0, imaging.Lanczos)
	imageBuffer := bytes.NewBuffer(nil)
	jpeg.Encode(imageBuffer, img, nil)
	return imageBuffer, nil
}

