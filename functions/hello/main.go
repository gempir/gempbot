package main

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func helloHandler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	id := os.Getenv("TWITCH_CLIENT_ID")

	return &events.APIGatewayProxyResponse{
		StatusCode:      200,
		Headers:         map[string]string{"Content-Type": "text/plain"},
		Body:            id,
		IsBase64Encoded: false,
	}, nil
}

func main() {
	lambda.Start(helloHandler)
}
