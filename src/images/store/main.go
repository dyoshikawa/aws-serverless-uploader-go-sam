package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dyoshikawa/aws-serverless-uploader-go/src/images/models"
	"github.com/dyoshikawa/aws-serverless-uploader-go/src/images/repositories"
	"github.com/dyoshikawa/aws-serverless-uploader-go/src/images/storage"
	"github.com/nfnt/resize"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/najeira/randstr"
)

const size = 300

type RequestBody struct {
	Data string `json:"data"`
}

type ResponseBody struct {
	Message string `json:"message"`
}

type ErrorResponseBody struct {
	Errors []string `json:"errors"`
}

type UploadInput struct {
	Base64File string
	Storage    storage.Storage
	Repo       repositories.Repository
	Region     string
	Bucket     string
}

func resizePng(file []byte) ([]byte, error) {
	img, err := png.Decode(bytes.NewReader(file))
	if err != nil {
		return nil, err
	}
	resizedImg := resize.Resize(size, 0, img, resize.Lanczos3)
	buf := new(bytes.Buffer)
	err = png.Encode(buf, resizedImg)
	if err != nil {
		return nil, err
	}
	resizedBytes := buf.Bytes()
	return resizedBytes, nil
}

func resizeJpeg(file []byte) ([]byte, error) {
	img, err := jpeg.Decode(bytes.NewReader(file))
	if err != nil {
		return nil, err
	}
	resizedImg := resize.Resize(size, 0, img, resize.Lanczos3)
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, resizedImg, nil)
	if err != nil {
		return nil, err
	}
	resizedBytes := buf.Bytes()
	return resizedBytes, nil
}

func upload(input UploadInput) error {
	strs := strings.Split(input.Base64File, ",")

	var ext string
	if strings.Index(strs[0], "image/png") != -1 {
		ext = ".png"
	} else if strings.Index(strs[0], "image/jpg") != -1 || strings.Index(strs[0], "image/jpeg") != -1 {
		ext = ".jpg"
	}
	file, err := base64.StdEncoding.DecodeString(strs[1])
	if err != nil {
		return err
	}
	name := randstr.String(50) + ext

	// 画像のリサイズ
	var resized []byte
	if ext == ".png" {
		resized, _ = resizePng(file)
	} else if ext == ".jpg" {
		resized, _ = resizeJpeg(file)
	}

	// S3アップロード
	err = input.Storage.Put(resized, name)
	if err != nil {
		return err
	}
	// DynamoDBにPut
	url := "https://s3-" + input.Region + ".amazonaws.com/" + input.Bucket + "/" + name
	input.Repo.Put(models.Image{
		Name:      name,
		URL:       url,
		CreatedAt: time.Now().Format("20060102150405"),
	})

	return nil
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
	fileStorage := os.Getenv("FILE_STORAGE_S3")
	if fileStorage == "" {
		msg := "Missed env variable FILE_STORAGE_S3."
		log.Println("ERROR: " + msg)
		return response(500, ErrorResponseBody{Errors: []string{"Missed settings."}}), errors.New(msg)
	}

	jsonBytes := ([]byte)(request.Body)
	reqBody := new(RequestBody)
	if err := json.Unmarshal(jsonBytes, reqBody); err != nil {
		log.Println("ERROR: " + err.Error())
		return response(500, ErrorResponseBody{
			Errors: []string{"ERROR: Failed JSON unmarshal."},
		}), err
	}

	// バリデーション
	strs := strings.Split(reqBody.Data, ",")
	if strings.Index(strs[0], "image/png") != -1 || strings.Index(strs[0], "image/jpg") != -1 || strings.Index(strs[0], "image/jpeg") != -1 {
	} else {
		return response(422, ErrorResponseBody{
			Errors: []string{"Given invalid file types."},
		}), nil
	}

	// Storage生成
	storage := storage.NewStorage(region, fileStorage)
	// Repository生成
	repo := repositories.NewRepository(region, appName)

	// メインロジック
	if err := upload(UploadInput{
		Base64File: reqBody.Data,
		Storage:    storage,
		Repo:       repo,
		Region:     region,
		Bucket:     fileStorage,
	}); err != nil {
		log.Println("ERROR: " + err.Error())
		return response(500, ErrorResponseBody{
			Errors: []string{"ERROR: Failed upload."},
		}), err
	}

	return response(200, ResponseBody{
		Message: "Success.",
	}), nil
}

func main() {
	lambda.Start(handler)
}
