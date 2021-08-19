package handlers

import (
	"aws-lambda-go/pkg/user"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"net/http"
)

var (
	ErrMethodNotAllow = errors.New("Method not allowed")
)

type ErrorBody struct {
	ErrorMsg string `json:"error,omitempty"`
}

func GetUser(req events.APIGatewayProxyRequest, tableName string, dClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	email := req.QueryStringParameters["email"]
	if len(email) > 0 {
		// Get a single user
		result, err := user.FetchUser(email, tableName, dClient)
		if err != nil {
			return apiResponse(http.StatusBadRequest, ErrorBody{*aws.String(err.Error())})
		}
		return apiResponse(http.StatusOK, result)
	}

	result, err := user.FetchUsers(tableName, dClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			*aws.String(err.Error()),
		})
	}
	return apiResponse(http.StatusOK, result)

}

func CreateUser(req events.APIGatewayProxyRequest, tableName string, dClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	result, err := user.CreateUser(req, tableName, dClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			*aws.String("Error when creating a user" + err.Error()),
		})
	}

	return apiResponse(http.StatusCreated, result)
}

func UpdateUser(req events.APIGatewayProxyRequest, tableName string, dClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	result, err := user.UpdateUser(req, tableName, dClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			*aws.String(err.Error()),
		})
	}

	return apiResponse(http.StatusOK, result)
}

func DeleteUser(req events.APIGatewayProxyRequest, tableName string, dClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	err := user.DeleteUser(req, tableName, dClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			*aws.String(err.Error()),
		})
	}

	return apiResponse(http.StatusOK, nil)
}

func UnhandledMethod() (*events.APIGatewayProxyResponse, error) {
	return apiResponse(http.StatusMethodNotAllowed, ErrMethodNotAllow)
}
