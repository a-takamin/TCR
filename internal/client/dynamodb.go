package client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func NewDynamoDbClient(isLocal bool) (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	return dynamodb.NewFromConfig(cfg, createDynamoDbOption(isLocal)), nil
}

func createDynamoDbOption(isLocal bool) func(o *dynamodb.Options) {
	if isLocal {
		return func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String("http://localhost:8000")
		}
	}
	return nil
}
