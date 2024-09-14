package repository

import (
	"context"
	"io"

	"github.com/a-takamin/tcr/internal/dto"
	"github.com/a-takamin/tcr/internal/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type BlobRepository struct {
	client     *s3.Client
	dClient    *dynamodb.Client
	bucketName string
	tableName  string
}

type Blob struct {
	Digest string
	Tag    string
	Blob   string
}

func NewBlobRepository(client *s3.Client, bucketName string, dynamodbClient *dynamodb.Client, tableName string) *BlobRepository {
	return &BlobRepository{
		client:     client,
		bucketName: bucketName,
		dClient:    dynamodbClient,
		tableName:  tableName,
	}
}

func (r BlobRepository) GetBlob(name string, digest string) (model.Blob, error) {
	resp, err := r.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(name + "/" + digest),
	})
	if err != nil {
		return model.Blob{}, err
	}
	blob, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Blob{}, err
	}

	return model.Blob{
		Name:   name,
		Digest: digest,
		Blob:   blob,
	}, nil
}

func (r BlobRepository) UploadBlob(key string, blob io.ReadCloser) error {
	_, err := r.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(key),
	})
	return err
}

func (r BlobRepository) GetChunkedBlobUploadProgress(name string) (dto.BlobUploadProgress, error) {
	resp, err := r.dClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"Name": &types.AttributeValueMemberS{
				Value: name,
			},
		},
	})
	if err != nil {
		return dto.BlobUploadProgress{}, err
	}

	var progress dto.BlobUploadProgress
	err = attributevalue.UnmarshalMap(resp.Item, &progress)
	if err != nil {
		return dto.BlobUploadProgress{}, err
	}
	return progress, nil
}

func (r BlobRepository) PutChunkedBlobUpdateProgress(newProgress dto.BlobUploadProgress) error {
	item, err := attributevalue.MarshalMap(newProgress)
	if err != nil {
		return err
	}
	_, err = r.dClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})

	return err
}

func (r BlobRepository) PutChunkedBlobConcatenateProgress(concatProgress dto.BlobConcatenateProgress) error {
	item, err := attributevalue.MarshalMap(concatProgress)
	if err != nil {
		return err
	}
	_, err = r.dClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("blob-concat-progress"),
		Item:      item,
	})
	return err
}

func (r BlobRepository) GetChunkedBlobConcatenateProgress(digest string) (dto.BlobConcatenateProgress, error) {
	resp, err := r.dClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("blob-concat-progress"),
		Key: map[string]types.AttributeValue{
			"Digest": &types.AttributeValueMemberS{
				Value: digest,
			},
		},
	})
	if err != nil {
		return dto.BlobConcatenateProgress{}, err
	}

	var progress dto.BlobConcatenateProgress
	err = attributevalue.UnmarshalMap(resp.Item, &progress)
	if err != nil {
		return dto.BlobConcatenateProgress{}, err
	}

	return progress, nil
}
