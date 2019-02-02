package repositories

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/dyoshikawa/aws-serverless-uploader-go/src/images/models"
	"github.com/guregu/dynamo"
)

type Repository interface {
	GetAll(data *models.Images) error
	Put(data models.Image) error
	DeleteAll(data []models.Image) error
}

type RepositoryDynamo struct {
	table dynamo.Table
}

func NewRepository(region string, appName string) Repository {
	db := dynamo.New(session.New(), &aws.Config{Region: aws.String(region)})

	repo := RepositoryDynamo{
		table: db.Table(appName + "Images"),
	}

	return &repo
}

func (repo *RepositoryDynamo) GetAll(data *models.Images) error {
	err := repo.table.Scan().All(data)
	if err != nil {
		return err
	}

	return nil
}

func (repo *RepositoryDynamo) Put(data models.Image) error {
	if err := repo.table.Put(data).Run(); err != nil {
		return err
	}

	return nil
}

func (repo *RepositoryDynamo) DeleteAll(data []models.Image) error {
	for _, v := range data {
		repo.table.Delete("Name", v.Name).Run()
	}

	return nil
}
