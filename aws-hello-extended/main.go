package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Age       int    `json:"age,omitempty"`
}

type Response struct {
	Name string `json:"name,omitempty"`
	Age  int    `json:"age,omitempty"`
}

func main() {

	lambda.Start(handle)
}

func handle(ctx context.Context, event Event) (string, error) {
	response := Response{
		Name: fmt.Sprintf("%s %s", event.FirstName, event.LastName),
		Age:  event.Age,
	}

	js, err := json.Marshal(response)
	if err != nil {
		return "", err
	}

	dataToString := string(js)
	return dataToString, nil
}
