package repository

import (
	"bytes"
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
	client                      *s3.Client
	dClient                     *dynamodb.Client
	bucketName                  string
	blobUploadProgressTableName string
	blobConcatProgressTableName string
}

type Blob struct {
	Digest string
	Tag    string
	Blob   string
}

func NewBlobRepository(client *s3.Client, bucketName string, dynamodbClient *dynamodb.Client, blobUploadProgressTName string, blobConcatProgressTName string) *BlobRepository {
	return &BlobRepository{
		client:                      client,
		bucketName:                  bucketName,
		dClient:                     dynamodbClient,
		blobUploadProgressTableName: blobUploadProgressTName,
		blobConcatProgressTableName: blobConcatProgressTName,
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

func (r BlobRepository) UploadBlob(key string, blob io.Reader) error {
	// TODO: なぜか分からないが「"failed to seek body to start, request stream is not seekable”」が発生するので、
	// 一度Blobを読み込んで再度Reader型にしている
	b, err := io.ReadAll(blob)
	if err != nil {
		return err
	}
	_, err = r.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(b),
	})
	return err
}

func (r BlobRepository) GetChunkedBlobUploadProgress(uuid string) (dto.BlobUploadProgress, error) {
	resp, err := r.dClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(r.blobUploadProgressTableName),
		Key: map[string]types.AttributeValue{
			"Uuid": &types.AttributeValueMemberS{
				Value: uuid,
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
		TableName: aws.String(r.blobUploadProgressTableName),
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
		TableName: aws.String(r.blobConcatProgressTableName),
		Item:      item,
	})
	return err
}

func (r BlobRepository) GetChunkedBlobConcatenateProgress(digest string) (dto.BlobConcatenateProgress, error) {
	resp, err := r.dClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(r.blobConcatProgressTableName),
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

func (r BlobRepository) DeleteBlob(input dto.DeleteBlobInput) error {
	_, err := r.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(input.Name + "/" + input.Digest),
	})
	return err
}
