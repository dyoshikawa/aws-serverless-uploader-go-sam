package storage

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Storage interface {
	Put(data []byte, key string) error
	Destroy() error
}

type StorageS3 struct {
	Svc            *s3.S3
	Iter           s3manager.BatchDeleteIterator
	PutObjectInput *s3.PutObjectInput
}

func NewStorage(region string, bucket string) Storage {
	svc := s3.New(session.New(), &aws.Config{
		Region: aws.String(region),
	})
	input := &s3.PutObjectInput{
		Bucket:               aws.String(bucket),
		ACL:                  aws.String("private"),
		ServerSideEncryption: aws.String("AES256"),
	}
	iter := s3manager.NewDeleteListIterator(svc, &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
	})

	return &StorageS3{
		Svc:            svc,
		Iter:           iter,
		PutObjectInput: input,
	}
}

func (storage *StorageS3) Put(data []byte, key string) error {
	input := storage.PutObjectInput
	input.Key = aws.String(key)
	input.Body = bytes.NewReader(data)
	_, err := storage.Svc.PutObject(input)
	if err != nil {
		return err
	}
	return nil
}

func (storage *StorageS3) Destroy() error {
	if err := s3manager.NewBatchDeleteWithClient(storage.Svc).Delete(aws.BackgroundContext(), storage.Iter); err != nil {
		return err
	}
	return nil
}
