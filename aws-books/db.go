package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var cfg = &aws.Config{
	Region: aws.String(endpoints.EuWest2RegionID),
}

// safe for concurrent use
var db = dynamodb.New(session.Must(session.NewSession(cfg)))

func getItem(isbn string) (*Book, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("Books"),
		Key: map[string]*dynamodb.AttributeValue{
			"ISBN": {S: aws.String(isbn)},
		},
	}

	result, err := db.GetItem(input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, nil
	}

	book := new(Book)

	err = dynamodbattribute.UnmarshalMap(result.Item, book)
	if err != nil {
		return nil, err
	}

	return book, nil

}

func putItemWithMarshalMap(bk *Book) error {
	data, err := dynamodbattribute.MarshalMap(bk)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("Books"),
		Item:      data,
	}

	_, err = db.PutItem(input)
	return err
}

//There are two options for putting the data into the database. Either we can allow aws to Marshal
//the data, in which case we called dynamodb.MarshalMap however we lose granuality. The major advantage
// of using Attribute value is that we set the Column names ourselves instead of relying on json names.
// I.e. authorName vs AuthorName. Also if we create custom types or want to use an enum it makes sense to
// construct the Item value(s) ourself.

func putItem(bk *Book) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String("Books"),
		Item: map[string]*dynamodb.AttributeValue{
			"ISBN": {
				S: aws.String(bk.ISBN),
			},
			"Title": {
				S: aws.String(bk.Title),
			},
			"Author": {
				S: aws.String(bk.Author)},
		},
	}
	_, err := db.PutItem(input)
	return err
}
