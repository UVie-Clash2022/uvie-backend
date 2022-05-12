package server

import (
	"github.com/aws/aws-lambda-go/events"
)

// Get500ServerError Internal server error
func Get500ServerError(err error) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: 500,
		Body:       err.Error(),
	}, nil
}

// Get400ServerError Bad request error
func Get400ServerError(msg string) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: 400,
		Body:       msg,
	}, nil
}
