package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
)

var (
	queueName = "sns-example-queue"
	topicARN  = "topic-arn"
)

func main() {
	cfg := &aws.Config{
		Region: aws.String(endpoints.EuWest2RegionID),
	}

	session := session.Must(session.NewSession(cfg))

	sqsClient := sqs.New(session)
	//snsClient := sns.New(session)

	queueURL, err := sqsClient.GetQueueUrl(&sqs.GetQueueUrlInput{QueueName: &queueName})
	if err != nil {
		panic(err)
	}

	sendMessageInput := &sqs.SendMessageInput{
		QueueUrl:    queueURL.QueueUrl,
		MessageBody: aws.String("Hey there litte buddy!!"),
	}

	_, err = sqsClient.SendMessage(sendMessageInput)

	if err != nil {
		panic(err)
	}

	receiveMessageInput := &sqs.ReceiveMessageInput{
		QueueUrl: queueURL.QueueUrl,
		AttributeNames: []*string{
			aws.String("All"),
		},
		MaxNumberOfMessages: aws.Int64(1),
		MessageAttributeNames: []*string{
			aws.String("All"),
		},
		VisibilityTimeout: aws.Int64(3600), // one hour to process request
	}

	receivedMessageOutput, err := sqsClient.ReceiveMessage(receiveMessageInput)

	if err != nil {
		panic(err)
	}

	log.Print(receivedMessageOutput)

}
