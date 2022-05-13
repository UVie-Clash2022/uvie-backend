package main

import (
	"context"
	"encoding/json"
	"github.com/UVie-Clash2022/uvie-backend/database"
	"github.com/UVie-Clash2022/uvie-backend/models"
	"github.com/UVie-Clash2022/uvie-backend/server"
	"github.com/aws/aws-lambda-go/events"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func ExcludeMovieForUser(request models.ExcludeMovieRequest) (*events.APIGatewayProxyResponse, error) {
	userCollection := database.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	if err := userCollection.FindOne(ctx, bson.M{"username": request.Username}).Decode(&user); err != nil {
		return server.Get400ServerError(err.Error())
	}

	validate := validator.New()
	validationErr := validate.Struct(&request)
	if validationErr != nil {
		return server.Get400ServerError(validationErr.Error())
	}

	//create map if it's null in the db
	if user.ExcludedMovies == nil {
		user.ExcludedMovies = map[string]struct{}{}
	}

	user.ExcludedMovies[request.MovieId] = struct{}{} //add movie to exclusion map

	_, err := userCollection.UpdateOne(ctx,
		bson.M{"_id": user.ID},
		bson.D{
			{"$set", bson.D{{"excludedMovies", user.ExcludedMovies}}},
		})

	if err != nil {
		return server.Get500ServerError(err)
	}

	return &events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
	}, nil
}

func GetExcludedMoviesForUser(username string) (*events.APIGatewayProxyResponse, error) {
	userCollection := database.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	if err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user); err != nil {
		return server.Get400ServerError(err.Error())
	}

	i := 0
	excludedMovieIds := make([]string, len(user.ExcludedMovies))
	for movieId := range user.ExcludedMovies {
		excludedMovieIds[i] = movieId
		i++
	}

	response := models.ExcludedMoviesResponse{Username: username, ExcludedMovieIds: excludedMovieIds}
	jsonResponse, err := json.Marshal(response)

	if err != nil {
		return server.Get500ServerError(err)
	}

	return &events.APIGatewayProxyResponse{
		StatusCode:      200,
		Headers:         map[string]string{"Content-Type": "application/json"},
		Body:            string(jsonResponse),
		IsBase64Encoded: false,
	}, nil
}
