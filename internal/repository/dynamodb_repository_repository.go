package repository

import (
	"context"

	"github.com/a-takamin/tcr/internal/dto"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Repository struct {
	Name string `dynamodbav:Name`
}

type RepositoryRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewRepositoryRepository(client *dynamodb.Client, TableName string) *RepositoryRepository {
	return &RepositoryRepository{
		client:    client,
		tableName: TableName,
	}
}

func (r RepositoryRepository) ExistsRepository(input dto.ExistsRepositoryInput) (bool, error) {
	itemInput := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"Name": &types.AttributeValueMemberS{
				Value: input.Name,
			},
		},
	}
	resp, err := r.client.GetItem(context.TODO(), itemInput)

	if err != nil {
		return false, err
	}

	var repo Repository
	err = attributevalue.UnmarshalMap(resp.Item, &repo)
	if err != nil {
		return false, err
	}

	if repo.Name == "" {
		return false, nil
	}

	return true, nil
}

func (r RepositoryRepository) SaveRepository(input dto.SaveRepositoryInput) error {
	repo := Repository{
		Name: input.Name,
	}
	item, err := attributevalue.MarshalMap(repo)
	if err != nil {
		return err
	}
	_, err = r.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	return err
}

func (r RepositoryRepository) DeleteRepository(input dto.DeleteRepositoryInput) error {
	itemInput := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"Name": &types.AttributeValueMemberS{
				Value: input.Name,
			},
		},
	}

	_, err := r.client.DeleteItem(context.TODO(), itemInput)
	return err
}
