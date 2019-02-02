package main

import (
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dyoshikawa/aws-serverless-uploader-go/src/images/models"
	"github.com/dyoshikawa/aws-serverless-uploader-go/src/images/repositories"
	"github.com/dyoshikawa/aws-serverless-uploader-go/src/images/storage"
)

func handler() error {
	appName := os.Getenv("APP_NAME")
	if appName == "" {
		msg := "Missed env variable APP_NAME."
		log.Fatalln("ERROR: " + msg)
	}
	region := os.Getenv("REGION")
	if region == "" {
		msg := "Missed env variable REGION."
		log.Fatalln("ERROR: " + msg)
	}
	fileStorage := os.Getenv("FILE_STORAGE_S3")
	if fileStorage == "" {
		msg := "Missed env variable FILE_STORAGE_S3."
		log.Fatalln("ERROR: " + msg)
	}

	storage := storage.NewStorage(region, fileStorage)
	if err := storage.Destroy(); err != nil {
		log.Fatalln("ERROR: " + err.Error())
	}

	repo := repositories.NewRepository(region, appName)
	var images models.Images
	repo.GetAll(&images)

	repo.DeleteAll(images)

	return nil
}

func main() {
	lambda.Start(handler)
}
