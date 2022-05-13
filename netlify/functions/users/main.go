package main

import (
	"encoding/json"
	"fmt"
	"github.com/UVie-Clash2022/uvie-backend/models"
	"github.com/UVie-Clash2022/uvie-backend/server"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"strings"
)

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	if request.HTTPMethod == "POST" {
		subPaths := strings.Split(request.Path, "/")
		method := subPaths[len(subPaths)-1]
		fmt.Println("method: " + method)

		var user models.User
		err := json.Unmarshal([]byte(request.Body), &user)
		if err != nil {
			return server.Get500ServerError(err)
		}

		if method == "login" {
			return LoginUser(user.Username, user.Password)
		} else if method == "signup" {
			return SignUpUser(user.Username, user.Password)
		} else {
			return server.Get400ServerError("Cannot process request")
		}
	} else if request.HTTPMethod == "PUT" {
		var excludeMovieRequest models.ExcludeMovieRequest
		err := json.Unmarshal([]byte(request.Body), &excludeMovieRequest)
		if err != nil {
			return server.Get500ServerError(err)
		}

		return ExcludeMovieForUser(excludeMovieRequest)
	} else if request.HTTPMethod == "GET" {
		subPaths := strings.Split(request.Path, "/")
		username := subPaths[len(subPaths)-2]
		fmt.Println("username: " + username)

		return GetExcludedMoviesForUser(username)
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 404,
	}, nil
}

func main() {
	lambda.Start(handler)
}
