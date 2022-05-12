package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"uvie-backend/server"
)

var MOVIEDB_API_KEY = os.Getenv("MOVIEDB_API_KEY")

func GetMovieData(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	movieQuery := request.QueryStringParameters["query"]

	if movieQuery == "" {
		return server.Get400ServerError("Movie query is empty")
	}

	searchQuery := url.QueryEscape(movieQuery)
	fullUrl := fmt.Sprintf("https://api.themoviedb.org/3/search/movie?api_key=%s&query=%s", MOVIEDB_API_KEY, searchQuery)

	response, err := http.Get(fullUrl)
	if err != nil {
		return server.Get500ServerError(err)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return server.Get500ServerError(err)
	}

	//formats json to be pretty-printed
	var out bytes.Buffer
	err = json.Indent(&out, responseData, "", "    ")
	if err != nil {
		return server.Get500ServerError(err)
	}

	return &events.APIGatewayProxyResponse{
		StatusCode:      200,
		Headers:         map[string]string{"Content-Type": "application/json"},
		Body:            out.String(),
		IsBase64Encoded: false,
	}, nil
}
