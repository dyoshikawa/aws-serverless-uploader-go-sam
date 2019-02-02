package main

import (
	"testing"

	"github.com/dyoshikawa/aws-serverless-uploader-go/src/images/repositories"
)

func TestIndex(t *testing.T) {
	t.Run("Able to do index", func(t *testing.T) {
		var repo repositories.Repository = &repositories.RepositoryMock{}

		_, err := index(repo)
		if err != nil {
			t.Fatal("Error: " + err.Error())
		}
	})
}
