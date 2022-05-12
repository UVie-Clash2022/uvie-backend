package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
)

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode:      200,
		Headers:         map[string]string{"Content-Type": "application/json"},
		Body:            fmt.Sprintf("{\"apiKey\":\"%s\"}", os.Getenv("MOVIEDB_API_KEY")),
		IsBase64Encoded: false,
	}, nil
}

func main() {
	lambda.Start(handler)
}
