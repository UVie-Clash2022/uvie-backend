package main

import (
	"encoding/json"
	"github.com/UVie-Clash2022/uvie-backend/models"
	"github.com/UVie-Clash2022/uvie-backend/server"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod == "POST" {
		var user models.User
		err := json.Unmarshal([]byte(request.Body), &user)
		if err != nil {
			return server.Get500ServerError(err)
		}

		return SignUpUser(user.Username, user.Password)
	}

	if request.HTTPMethod == "GET" {
		username := request.QueryStringParameters["username"]
		if username == "" {
			return server.Get400ServerError("username is empty")
		}

		return GetUser(username)
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 404,
	}, nil
}

func main() {
	lambda.Start(handler)
}
