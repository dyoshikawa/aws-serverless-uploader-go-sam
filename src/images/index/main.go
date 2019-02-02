package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sort"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dyoshikawa/aws-serverless-uploader-go/src/images/models"
	"github.com/dyoshikawa/aws-serverless-uploader-go/src/images/repositories"
)

type ErrorResponseBody struct {
	Errors []string `json:"errors"`
}

func index(repo repositories.Repository) ([]models.Image, error) {
	var images models.Images
	err := repo.GetAll(&images)
	if err != nil {
		return []models.Image{}, err
	}

	sort.Sort(images)

	return images, nil
}

func response(code int, body interface{}) events.APIGatewayProxyResponse {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		log.Fatalln("ERROR: " + err.Error())
	}

	return events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "origin,Accept,Authorization,Content-Type",
			"Content-Type":                 "application/json",
		},
		Body:       string(jsonBytes),
		StatusCode: code,
	}
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	appName := os.Getenv("APP_NAME")
	if appName == "" {
		msg := "Missed env variable APP_NAME."
		log.Println("ERROR: " + msg)
		return response(500, ErrorResponseBody{Errors: []string{"Missed settings."}}), errors.New(msg)
	}
	region := os.Getenv("REGION")
	if region == "" {
		msg := "Missed env variable REGION."
		log.Println("ERROR: " + msg)
		return response(500, ErrorResponseBody{Errors: []string{"Missed settings."}}), errors.New(msg)
	}

	// Repository生成
	repo := repositories.NewRepository(region, appName)

	// メインロジック
	images, err := index(repo)
	if err != nil {
		log.Println("ERROR: " + err.Error())
		return response(500, ErrorResponseBody{Errors: []string{"Failed to index."}}), nil
	}

	return response(200, images), nil
}

func main() {
	lambda.Start(handler)
}
