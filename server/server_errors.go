package server

import (
	"encoding/json"
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
func Get400ServerError(errorMessage string) (*events.APIGatewayProxyResponse, error) {
	responseMap := map[string]string{
		"errorMessage": errorMessage,
	}
	jsonResponse, _ := json.Marshal(responseMap)

	return &events.APIGatewayProxyResponse{
		StatusCode: 400,
		Body:       string(jsonResponse),
	}, nil
}
