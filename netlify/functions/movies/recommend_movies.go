package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/UVie-Clash2022/uvie-backend/database"
	"github.com/UVie-Clash2022/uvie-backend/models"
	"github.com/UVie-Clash2022/uvie-backend/server"
	"github.com/aws/aws-lambda-go/events"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func RecommendMovies(username string, page int) (*events.APIGatewayProxyResponse, error) {
	//get user's liked movies
	moviesUserCollection := database.GetCollection("movies_user")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var userMoviesLiked models.UserMoviesLiked
	if err := moviesUserCollection.FindOne(ctx, bson.M{"username": username}).Decode(&userMoviesLiked); err != nil {
		return server.Get400ServerError(err.Error())
	}

	//get user's disliked or "excluded" move list
	userCollection := database.GetCollection("users")
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	if err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user); err != nil {
		return server.Get400ServerError(err.Error())
	}
	dislikedMovieIds := user.ExcludedMovies

	//get user's current favorite movie to base recommendations on
	currentFavMovie := userMoviesLiked.CurrentFavorite

	var recommendMoviesDbResponse models.RecommendMoviesResponse
	userRecommendationsResponse := models.RecommendMoviesResponse{}
	goToNextPage := true

	for goToNextPage {
		//Call themoviedb api to get recommended movies based on the user's favorite movieId
		recommendEndpoint := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d/recommendations?api_key=%s&language=en-US&page=%d",
			currentFavMovie.ID, MOVIEDB_API_KEY, page)

		response, err := http.Get(recommendEndpoint)
		if err != nil {
			return server.Get500ServerError(err)
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return server.Get500ServerError(err)
		}

		err = json.Unmarshal(responseData, &recommendMoviesDbResponse)
		if err != nil {
			return server.Get500ServerError(err)
		}

		//Populate recommendation list and check if each recommendation is already liked or disliked, if so, exclude it from list.
		//If all page results have been iterated and there are still less than 10 recommendations, get the next page of results.
		goToNextPage = buildRecommendationList(&recommendMoviesDbResponse, &userRecommendationsResponse, userMoviesLiked.Liked, dislikedMovieIds, page)
		if goToNextPage {
			page++
		}
	}
	userRecommendationsResponse.Page = page
	userRecommendationsResponse.TotalPages = recommendMoviesDbResponse.TotalPages
	userRecommendationsResponse.TotalResults = recommendMoviesDbResponse.TotalResults

	jsonResponse, err := json.Marshal(userRecommendationsResponse)
	if err != nil {
		return server.Get500ServerError(err)
	}

	//return list of 10 recommendations
	return &events.APIGatewayProxyResponse{
		StatusCode:      200,
		Headers:         map[string]string{"Content-Type": "application/json"},
		Body:            string(jsonResponse),
		IsBase64Encoded: false,
	}, nil
}

func buildRecommendationList(recDbResponse *models.RecommendMoviesResponse,
	finalResponse *models.RecommendMoviesResponse,
	likedMovies []models.Movie,
	dislikedMovies map[string]struct{},
	page int) bool {

	for _, movie := range recDbResponse.Movies {
		if !existsInLikedMovies(likedMovies, movie.ID) && !existsInDislikedMovies(dislikedMovies, movie.ID) {
			movie.PosterUrl = "https://image.tmdb.org/t/p/original" + movie.PosterUrl
			finalResponse.Movies = append(finalResponse.Movies, movie)
			if len(finalResponse.Movies) >= 10 {
				return false //goToNextPage is set to false because we got 10 new movies in the list
			}
		}
	}
	if recDbResponse.TotalPages > page { //we did not build a list with 10 movies so check if there are more pages
		return true //goToNextPage is set to true
	}

	return false //no more pages to process, goToNextPage is set to false
}

func existsInLikedMovies(likedMovies []models.Movie, movieId int) bool {
	for _, v := range likedMovies {
		if movieId == v.ID {
			return true
		}
	}
	return false
}

func existsInDislikedMovies(dislikedMovies map[string]struct{}, movieId int) bool {
	if _, ok := dislikedMovies[strconv.Itoa(movieId)]; ok {
		return true
	}
	return false
}
