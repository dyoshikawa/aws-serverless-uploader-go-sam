package repositories

import "github.com/dyoshikawa/aws-serverless-uploader-go/src/images/models"

type RepositoryMock struct{}

func (repo *RepositoryMock) GetAll(data *models.Images) error {
	return nil
}

func (repo *RepositoryMock) Put(data models.Image) error {
	return nil
}
