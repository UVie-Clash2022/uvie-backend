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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func LoginUser(username string, password string) (*events.APIGatewayProxyResponse, error) {
	userCollection := database.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	if err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user); err != nil {
		return server.Get400ServerError(err.Error())
	}

	if user.Password != password {
		return server.Get400ServerError("Invalid username or password")
	}

	jsonUser, err := json.Marshal(user)
	if err != nil {
		return server.Get500ServerError(err)
	}

	return &events.APIGatewayProxyResponse{
		StatusCode:      200,
		Headers:         map[string]string{"Content-Type": "application/json"},
		Body:            string(jsonUser),
		IsBase64Encoded: false,
	}, nil
}

func SignUpUser(username string, password string) (*events.APIGatewayProxyResponse, error) {
	userCollection := database.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userCount, err := userCollection.CountDocuments(ctx, bson.M{"username": username})
	if err != nil {
		return server.Get500ServerError(err)
	}

	if userCount > 0 {
		return server.Get400ServerError(fmt.Sprintf("Cannot sign up because the username '%s' already exists.", username))
	}

	user := models.User{
		ID:       primitive.NewObjectID(),
		Username: username,
		Password: password,
	}

	validate := validator.New()
	validationErr := validate.Struct(&user)
	if validationErr != nil {
		return server.Get400ServerError(validationErr.Error())
	}

	_, err = userCollection.InsertOne(ctx, user)
	if err != nil {
		return server.Get500ServerError(err)
	}

	responseMap := map[string]string{
		"msg": "Successfully signed up!",
	}
	jsonResponse, _ := json.Marshal(responseMap)

	return &events.APIGatewayProxyResponse{
		StatusCode:      200,
		Headers:         map[string]string{"Content-Type": "application/json"},
		Body:            string(jsonResponse),
		IsBase64Encoded: false,
	}, nil
}
