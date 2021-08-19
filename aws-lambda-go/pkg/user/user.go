package user

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"regexp"
)

var (
	ErrorFailedToUnmarshalRecord = "failed to unmarshal record"
	ErrorFailedToFetchRecord     = "failed to fetch record"
	ErrorInvalidUserData         = "invalid user data"
	ErrorInvalidEmail            = "invalid email"
	ErrorCouldNotMarshalItem     = "could not marshal item"
	ErrorCouldNotDeleteItem      = "could not delete item"
	ErrorCouldNotDynamoPutItem   = "could not dynamo put item error"
	ErrorUserAlreadyExists       = "user.User already exists"
	ErrorUserDoesNotExists       = "user.User does not exist"
)

type User struct {
	Email     string `json:"email,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}

func FetchUser(email, tableName string, dClient dynamodbiface.DynamoDBAPI) (*User, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {S: aws.String(email)},
		},
		TableName: aws.String(tableName),
	}

	result, err := dClient.GetItem(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	user := new(User)
	// The output value provided must be a non-nil pointer!!!!!
	// If you use user := User{} it will fail but not generate an error
	err = dynamodbattribute.UnmarshalMap(result.Item, user)

	return user, nil
}

func FetchUsers(tableName string, dClient dynamodbiface.DynamoDBAPI) (*[]User, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}
	result, err := dClient.Scan(input)

	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	// Pass in a NON NIL POINTER to the umarshal map!!!
	users := new([]User)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, users)
	if err != nil {
		return nil, errors.New("Failed when unmarshalling list of users")
	}

	return users, nil

}

func CreateUser(req events.APIGatewayProxyRequest, tableName string, dClient dynamodbiface.DynamoDBAPI) (*User, error) {
	/*
		Create an Empty User then unmarshal req body to user{}, validating their email
		Check if the User already exists by checking DB for their email address
		Marshal the User using dynamodbattribute.MarhsalMap
		Create a dynamodb PutItem with the MarhsalMap Item and the TableName
		Use the dynamodbclient to Put this item ^
	*/

	// This code is not nice. Why should the user function unmarshal the incoming request!
	// This is the job of the handler.
	var user User
	// note field body is a string, hence string -> byte[]
	if err := json.Unmarshal([]byte(req.Body), &user); err != nil {
		return nil, errors.New(ErrorInvalidUserData)
	}

	if !IsEmailValid(user.Email) {
		return nil, errors.New(ErrorInvalidEmail)
	}

	// does the current user exist?
	currentUser, err := FetchUser(user.Email, tableName, dClient)
	if currentUser != nil && len(currentUser.Email) != 0 {
		return nil, errors.New(ErrorUserAlreadyExists)
	}

	d, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	input := dynamodb.PutItemInput{
		Item:      d,
		TableName: aws.String(tableName),
	}

	_, err = dClient.PutItem(&input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}

	return &user, nil

}

func UpdateUser(req events.APIGatewayProxyRequest, tableName string, dClient dynamodbiface.DynamoDBAPI) (*User, error) {
	user := User{}
	if err := json.Unmarshal([]byte(req.Body), &user); err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	// does the user exist?
	currentUser, _ := FetchUser(user.Email, tableName, dClient)
	if currentUser != nil && len(currentUser.Email) == 0 {
		return nil, errors.New(ErrorUserDoesNotExists)
	}

	d, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	input := dynamodb.PutItemInput{
		Item:      d,
		TableName: aws.String(tableName),
	}

	_, err = dClient.PutItem(&input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}

	return &user, nil
}

func DeleteUser(req events.APIGatewayProxyRequest, tableName string, dClient dynamodbiface.DynamoDBAPI) error {
	email := req.QueryStringParameters["email"]
	input := dynamodb.DeleteItemInput{Key: map[string]*dynamodb.AttributeValue{
		"email": {S: aws.String(email)},
	},
		TableName: aws.String(tableName)}

	_, err := dClient.DeleteItem(&input)
	if err != nil {
		return errors.New(ErrorCouldNotDeleteItem)
	}
	return nil
}

func IsEmailValid(email string) bool {
	var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]{1,64}@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if len(email) < 3 || len(email) > 254 || !rxEmail.MatchString(email) {
		return false
	}
	return true
}
