package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// NotificationRequest is the expected structure of the incoming request body.
type NotificationRequest struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// SNSClient is an interface for the SNS Publish operation, for testability.
type SNSClient interface {
	Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}

var snsClient SNSClient

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// 1. Get SNS Topic ARN from environment variables
	topicArn := os.Getenv("SNS_TOPIC_ARN")
	if topicArn == "" {
		return events.APIGatewayProxyResponse{StatusCode: 500}, fmt.Errorf("SNS_TOPIC_ARN environment variable not set")
	}

	// 2. Parse the incoming request
	var req NotificationRequest
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("failed to parse request body: %w", err)
	}

	if req.Title == "" || req.Message == "" {
		return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("title and message are required")
	}

	// 3. Publish to SNS
	_, err = snsClient.Publish(ctx, &sns.PublishInput{
		Message:  &req.Message,
		Subject:  &req.Title,
		TopicArn: &topicArn,
	})

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, fmt.Errorf("failed to publish to SNS: %w", err)
	}

	// 4. Return a success response
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       `{"status":"notification published"}`,
	}, nil
}

func main() {
	// Initialize the SNS client once
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(fmt.Sprintf("unable to load SDK config, %v", err))
	}
	snsClient = sns.NewFromConfig(cfg)

	lambda.Start(handler)
}
