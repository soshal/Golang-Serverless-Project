package main

import (
	"log"
	"os"

	"github.com/akhil/go-serverless-yt/pkg/handlers"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	dynaClient dynamodbiface.DynamoDBAPI
)

func main() {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		log.Fatal("AWS_REGION environment variable not set")
	}

	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		log.Fatalf("Failed to create AWS session: %s", err)
	}

	dynaClient = dynamodb.New(awsSession)
	lambda.Start(handler)
}

const tableName = "go-serverless-yt"

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	methodHandlers := map[string]func(events.APIGatewayProxyRequest, string, dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error){
		"GET":    handlers.GetUser,
		"POST":   handlers.CreateUser,
		"PUT":    handlers.UpdateUser,
		"DELETE": handlers.DeleteUser,
	}

	if handlerFunc, exists := methodHandlers[req.HTTPMethod]; exists {
		return handlerFunc(req, tableName, dynaClient)
	}
	return handlers.UnhandledMethod()
}
