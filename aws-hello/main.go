package main

import (
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"net/http"
)

var (
	ErrNameNotProvided = errors.New("no name supplied to the request body")
)

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Processing your lambda request %s\n", request.RequestContext.RequestID)

	if len(request.Body) < 1 {
		return events.APIGatewayProxyResponse{}, ErrNameNotProvided
	}

	response := events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Hello %s", request.Body),
		StatusCode: http.StatusOK,
	}

	return response, nil
}
