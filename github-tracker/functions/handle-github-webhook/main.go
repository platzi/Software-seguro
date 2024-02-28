package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		IsBase64Encoded: false,
		StatusCode:      http.StatusOK,
		Body:            "hello from handle-github-notifications",
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
