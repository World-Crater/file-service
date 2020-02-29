package image

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
    "github.com/google/uuid"
    "fmt"
    "errors"
)

type Image struct {
    ID     string `dynamodbav:"id"`
    Url    string `dynamodbav:"url"`
}

func SaveImage(tableName string, url string) (*dynamodb.PutItemOutput, error){
    id, err := uuid.NewUUID()
    if err != nil {
        fmt.Println("Got uuid error")
    }
    // snippet-start:[dynamodb.go.create_item.session]
    // Initialize a session that the SDK will use to load
    // credentials from the shared credentials file ~/.aws/credentials
    // and region from the shared configuration file ~/.aws/config.
    sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))

    // Create DynamoDB client
    svc := dynamodb.New(sess)
    // snippet-end:[dynamodb.go.create_item.session]

    // snippet-start:[dynamodb.go.create_item.assign_struct]
    item := Image{
        ID:     id.String(),
        Url:    url,
    }

    av, err := dynamodbattribute.MarshalMap(item)
    if err != nil {
        fmt.Println("Got error marshalling new movie item:")
        fmt.Println(err.Error())
        return nil, errors.New("Got error marshalling new movie item:")
    }
    fmt.Println("%d", av)

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