package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod == "GET" {
		return GetMovieData(request)
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 404,
	}, nil
}

func main() {
	lambda.Start(handler)
}
