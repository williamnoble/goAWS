package handlers

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"log"
)

func apiResponse(status int, body interface{}) (*events.APIGatewayProxyResponse, error) {

	d, err := json.Marshal(body)
	if err != nil {
		log.Println("error in marshalling response body")
	}

	resp := events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StatusCode: status,
		// AWS uses body string so convert from []byte
		Body: string(d),
	}

	return &resp, nil
}
