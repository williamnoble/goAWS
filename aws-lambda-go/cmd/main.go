package main

import (
	"aws-lambda-go/pkg/handlers"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

const (
	tableName = "LambdaGoUser"
)

// This is ugly.
var (
	dClient dynamodbiface.DynamoDBAPI
)

func main() {

	var cfg = &aws.Config{
		Region: aws.String(endpoints.EuWest2RegionID),
	}

	dClient = dynamodb.New(session.Must(session.NewSession(cfg)))
	lambda.Start(router)
}

func router(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return handlers.GetUser(req, tableName, dClient)
	case "POST":
		return handlers.CreateUser(req, tableName, dClient)
	case "PUT":
		return handlers.UpdateUser(req, tableName, dClient)
	case "DELETE":
		return handlers.DeleteUser(req, tableName, dClient)
	default:
		return handlers.UnhandledMethod()
	}
}
