package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/UVie-Clash2022/uvie-backend/database"
	"github.com/UVie-Clash2022/uvie-backend/models"
	"github.com/UVie-Clash2022/uvie-backend/server"
	"github.com/aws/aws-lambda-go/events"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func GetMoviesLikedForUser(username string) (*events.APIGatewayProxyResponse, error) {
	moviesUserCollection := database.GetCollection("movies_user")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var userMoviesLiked models.UserMoviesLiked
	if err := moviesUserCollection.FindOne(ctx, bson.M{"username": username}).Decode(&userMoviesLiked); err != nil {
		return server.Get400ServerError(err.Error())
	}

	userMoviesLiked.Username = "" //so it's removed from response
	jsonResponse, err := json.Marshal(userMoviesLiked)
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

func SaveLikedMovieForUser(request models.LikedMovieRequest) (*events.APIGatewayProxyResponse, error) {
	validate := validator.New()
	validationErr := validate.Struct(&request)
	if validationErr != nil {
		return server.Get400ServerError(validationErr.Error())
	}

	moviesUserCollection := database.GetCollection("movies_user")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userMoviesLikesDocExists := true
	var userMoviesLiked models.UserMoviesLiked
	if err := moviesUserCollection.FindOne(ctx, bson.M{"username": request.Username}).Decode(&userMoviesLiked); err != nil {
		if err == mongo.ErrNoDocuments {
			userMoviesLikesDocExists = false
			userMoviesLiked = models.UserMoviesLiked{Username: request.Username}
		} else {
			return server.Get400ServerError(err.Error())
		}
	}

	// Get movie data by using the movieId from request
	movieDbEndpoint := fmt.Sprintf("https://api.themoviedb.org/3/movie/%s?api_key=%s&language=en-US", request.MovieId, MOVIEDB_API_KEY)
	response, err := http.Get(movieDbEndpoint)
	if err != nil {
		return server.Get500ServerError(err)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return server.Get500ServerError(err)
	}

	var likedMovie models.Movie
	err = json.Unmarshal(responseData, &likedMovie)
	if err != nil {
		return server.Get500ServerError(err)
	}

	likedMovie.PosterUrl = "https://image.tmdb.org/t/p/original" + likedMovie.PosterUrl
	//fmt.Printf("%+v\n", likedMovie)

	//look for duplicated likes
	duplicated := false
	for _, v := range userMoviesLiked.Liked {
		if v.ID == likedMovie.ID {
			duplicated = true
			break
		}
	}
	//Append db liked array with the new liked movie
	if !duplicated {
		userMoviesLiked.Liked = append(userMoviesLiked.Liked, likedMovie)
	}

	//Update the currentFavorite movie in the db item
	userMoviesLiked.CurrentFavorite = likedMovie

	if userMoviesLikesDocExists {
		//save changes to mongo collection
		_, err = moviesUserCollection.UpdateOne(
			ctx,
			bson.M{"username": userMoviesLiked.Username},
			bson.D{
				{"$set", bson.D{
					{"currentFavorite", userMoviesLiked.CurrentFavorite},
					{"liked", userMoviesLiked.Liked},
				}},
			})
	} else {
		_, err = moviesUserCollection.InsertOne(ctx, userMoviesLiked)
	}

	if err != nil {
		return server.Get500ServerError(err)
	}

	return &events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
	}, nil
}

func RemoveLikedMovieForUser(request models.LikedMovieRequest) (*events.APIGatewayProxyResponse, error) {
	validate := validator.New()
	validationErr := validate.Struct(&request)
	if validationErr != nil {
		return server.Get400ServerError(validationErr.Error())
	}

	movieIdToDelete, err := strconv.Atoi(request.MovieId)
	if err != nil {
		return server.Get400ServerError(err.Error())
	}

	moviesUserCollection := database.GetCollection("movies_user")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var userMoviesLiked models.UserMoviesLiked
	if err := moviesUserCollection.FindOne(ctx, bson.M{"username": request.Username}).Decode(&userMoviesLiked); err != nil {
		return server.Get400ServerError(err.Error())
	}

	//Remove liked movie from list
	idxToRemove := 0
	for _, likedMovie := range userMoviesLiked.Liked {
		if likedMovie.ID == movieIdToDelete {
			break
		}
		idxToRemove++
	}
	if idxToRemove >= len(userMoviesLiked.Liked) {
		return server.Get400ServerError(fmt.Sprintf("The provided movieId was not found in the user's liked list."))
	}

	userMoviesLiked.Liked = append(userMoviesLiked.Liked[:idxToRemove], userMoviesLiked.Liked[idxToRemove+1:]...)

	//todo: should we do below: update currentFavorite movie if the removed liked movie is the currentFavorite?
	if userMoviesLiked.CurrentFavorite.ID == movieIdToDelete {
		userMoviesLiked.CurrentFavorite = userMoviesLiked.Liked[len(userMoviesLiked.Liked)-1]
	}

	//update db with new list
	_, err = moviesUserCollection.UpdateOne(
		ctx,
		bson.M{"username": request.Username},
		bson.D{
			{"$set", bson.D{
				{"currentFavorite", userMoviesLiked.CurrentFavorite},
				{"liked", userMoviesLiked.Liked},
			}},
		})

	if err != nil {
		return server.Get500ServerError(err)
	}

	return &events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
	}, nil
}
