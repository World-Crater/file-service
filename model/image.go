package image

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type Image struct {
	ID  string `dynamodbav:"id"`
	Url string `dynamodbav:"url"`
}

func SaveImage(tableName string, url string) (*dynamodb.PutItemOutput, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		fmt.Println("Got uuid error")
	}
	// snippet-start:[dynamodb.go.create_item.session]
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(viper.GetString("AWS_REGION")),
	})

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	// snippet-end:[dynamodb.go.create_item.session]

	// snippet-start:[dynamodb.go.create_item.assign_struct]
	item := Image{
		ID:  id.String(),
		Url: url,
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		fmt.Println("Got error marshalling new movie item:")
		fmt.Println(err.Error())
		return nil, errors.New("Got error marshalling new movie item:")
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	putResult, err := svc.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		return nil, errors.New("Got error calling PutItem:")
	}

	fmt.Println("Successfully added")
	return putResult, nil
}

func SaveImageToS3(image *bytes.Buffer, key string) (*s3manager.UploadOutput, error) {
	// The session the S3 Uploader will use
	sess := session.Must(session.NewSession())

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)
	readBuf, _ := ioutil.ReadAll(image)
	f := bytes.NewReader(readBuf)

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(viper.GetString("IMAGE_BUCKET_NAME")),
		Key:    aws.String(key),
		Body:   f,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file, %v", err)
	}
	return result, nil
}
