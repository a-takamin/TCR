package repository

import (
	"context"

	"github.com/a-takamin/tcr/internal/dto"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type BlobUploadProgress struct {
	Uuid         string `dynamodbav:"Uuid"`
	ByteUploaded int64  `dynamodbav:"ByteUploaded"`
	NextChunkNo  int    `dynamodbav:"NextChunkNo"`
	Done         bool   `dynamodbav:"Done"`
	Digest       string `dynamodbav:"Digest"`
}

type BlobUploadProgressRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewBlobUploadProgressRepository(client *dynamodb.Client, TableName string) *BlobUploadProgressRepository {
	return &BlobUploadProgressRepository{
		client:    client,
		tableName: TableName,
	}
}

func (r BlobUploadProgressRepository) FindBlobUploadProgress(input dto.FindBlobUploadProgressInput) (dto.FindBlobUploadProgressOutput, error) {
	resp, err := r.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"Uuid": &types.AttributeValueMemberS{
				Value: input.Uuid,
			},
		},
	})
	if err != nil {
		return dto.FindBlobUploadProgressOutput{}, err
	}

	var progress dto.BlobUploadProgress
	err = attributevalue.UnmarshalMap(resp.Item, &progress)
	if err != nil {
		return dto.FindBlobUploadProgressOutput{}, err
	}

	return dto.FindBlobUploadProgressOutput{
		Uuid:         progress.Uuid,
		ByteUploaded: progress.ByteUploaded,
		NextChunkNo:  progress.NextChunkNo,
		Digest:       progress.Digest,
	}, nil
}

func (r BlobUploadProgressRepository) SaveBlobUploadProgress(input dto.SaveBlobUploadProgressInput) error {
	progress := BlobUploadProgress{
		Uuid:         input.Uuid,
		ByteUploaded: input.ByteUploaded,
		NextChunkNo:  input.NextChunkNo,
		Digest:       input.Digest,
	}
	item, err := attributevalue.MarshalMap(progress)
	if err != nil {
		return err
	}
	_, err = r.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})

	return err
}
