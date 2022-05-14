package main

import (
	"encoding/json"
	"github.com/UVie-Clash2022/uvie-backend/models"
	"github.com/UVie-Clash2022/uvie-backend/server"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"strconv"
	"strings"
)

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod == "GET" {
		subPaths := strings.Split(request.Path, "/")
		method := subPaths[len(subPaths)-2]

		if method == "movie-likes" {
			username := subPaths[len(subPaths)-1]
			return GetMoviesLikedForUser(username)
		} else if method == "recommend" {
			username := subPaths[len(subPaths)-1]
			pageQueryParam := request.QueryStringParameters["page"]
			if pageQueryParam == "" {
				return RecommendMovies(username, 1)
			} else {
				page, err := strconv.Atoi(pageQueryParam)
				if err != nil {
					return server.Get400ServerError(err.Error())
				}
				return RecommendMovies(username, page)
			}
		}

		return GetMovieData(request)
	} else if request.HTTPMethod == "PUT" {
		var movieLikedRequest models.LikedMovieRequest
		err := json.Unmarshal([]byte(request.Body), &movieLikedRequest)
		if err != nil {
			return server.Get500ServerError(err)
		}
		return SaveLikedMovieForUser(movieLikedRequest)
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 404,
	}, nil
}

func main() {
	lambda.Start(handler)
}
