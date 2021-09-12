package main

import (
	"github.com/aws/aws-lambda-go/events"
)

func worldHandler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode:      200,
		Headers:         map[string]string{"Content-Type": "text/plain"},
		Body:            "Hello, World!",
		IsBase64Encoded: false,
	}, nil
}
