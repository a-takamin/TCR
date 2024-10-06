package repository

import (
	"context"
	"encoding/base64"

	"github.com/a-takamin/tcr/internal/dto"
	"github.com/a-takamin/tcr/internal/service/domain"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type ManifestRepository struct {
	client            *dynamodb.Client
	manifestTableName string
}

type Manifest struct {
	Name     string `dynamodbav:Name`
	Digest   string `dynamodbav:Digest`
	Tag      string `dynamodbav:Tag`
	Manifest string `dynamodbav:Manifest`
}

func NewManifestRepository(client *dynamodb.Client, manifestTableName string) *ManifestRepository {
	return &ManifestRepository{
		client:            client,
		manifestTableName: manifestTableName,
	}
}

func (r ManifestRepository) getItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	return r.client.GetItem(ctx, params, optFns...)
}

func (r ManifestRepository) QueryItem(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	return r.client.Query(ctx, params, optFns...)
}

func (r ManifestRepository) GetTags(name string) (dto.GetTagsResponse, error) {
	input, err := r.createGetTagsInput(name)
	if err != nil {
		return dto.GetTagsResponse{}, err
	}
	resp, err := r.QueryItem(context.TODO(), input)
	if err != nil {
		return dto.GetTagsResponse{}, err
	}
	var manifests []Manifest
	err = attributevalue.UnmarshalListOfMaps(resp.Items, &manifests)
	if err != nil {
		return dto.GetTagsResponse{}, err
	}

	var tags dto.GetTagsResponse
	for _, m := range manifests {
		tags.Tags = append(tags.Tags, m.Tag)
	}

	return tags, nil
}

func (r ManifestRepository) createGetTagsInput(name string) (*dynamodb.QueryInput, error) {
	keyEx := expression.Key("Name").Equal(expression.Value(name))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return &dynamodb.QueryInput{}, err
	}
	return &dynamodb.QueryInput{
		TableName:                 aws.String(r.manifestTableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}, nil
}

// リファクタ
// / Name があるかどうかを確認する関数
// リファクタメモここまで

func (r ManifestRepository) ExistsName(name string) (bool, error) {
	keyEx := expression.Key("Name").Equal(expression.Value(name))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return false, err
	}
	input := &dynamodb.QueryInput{
		TableName:                 aws.String(r.manifestTableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}
	resp, err := r.QueryItem(context.TODO(), input)
	if err != nil {
		return false, err
	}
	if resp.Count == 0 {
		return false, nil
	}
	return true, nil
}

// リファクタ
func (r ManifestRepository) ExistsManifest(input dto.ExistsManifestInput) (bool, error) {
	manifest, err := r.FindManifest(dto.FindManifestInput{
		Name:      input.Name,
		Reference: input.Reference,
	})
	if err != nil {
		return false, err
	}
	if manifest.Name == "" {
		return false, nil
	}
	return true, nil
}

func (r ManifestRepository) FindManifest(input dto.FindManifestInput) (dto.FindManifestOutput, error) {
	if domain.IsDigest(input.Reference) {
		return r.FindManifestByDigest(input)
	} else {
		return r.FindManifestByTag(input)
	}
}

func (r ManifestRepository) FindManifestByDigest(input dto.FindManifestInput) (dto.FindManifestOutput, error) {
	itemInput := &dynamodb.GetItemInput{
		TableName: aws.String(r.manifestTableName),
		Key: map[string]types.AttributeValue{
			"Name": &types.AttributeValueMemberS{
				Value: input.Name,
			},
			"Digest": &types.AttributeValueMemberS{
				Value: input.Reference,
			},
		},
	}

	resp, err := r.getItem(context.TODO(), itemInput)

	if err != nil {
		return dto.FindManifestOutput{}, err
	}

	var dbManifest Manifest
	err = attributevalue.UnmarshalMap(resp.Item, &dbManifest)
	if err != nil {
		return dto.FindManifestOutput{}, err
	}

	decordedManifest, err := base64.StdEncoding.DecodeString(dbManifest.Manifest)
	if err != nil {
		return dto.FindManifestOutput{}, err
	}

	return dto.FindManifestOutput{
		Name:     dbManifest.Name,
		Tag:      dbManifest.Tag,
		Digest:   dbManifest.Digest,
		Manifest: decordedManifest,
	}, nil
}

func (r ManifestRepository) FindManifestByTag(input dto.FindManifestInput) (dto.FindManifestOutput, error) {
	keyEx := expression.KeyAnd(
		expression.Key("Name").Equal(expression.Value(input.Name)),
		expression.Key("Tag").Equal(expression.Value(input.Reference)),
	)
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return dto.FindManifestOutput{}, err
	}
	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(r.manifestTableName),
		IndexName:                 aws.String("ManifestTagIndex"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}
	resp, err := r.QueryItem(context.TODO(), queryInput)

	if err != nil {
		return dto.FindManifestOutput{}, err
	}
	var manifests []Manifest
	err = attributevalue.UnmarshalListOfMaps(resp.Items, &manifests)
	if err != nil {
		return dto.FindManifestOutput{}, err
	}

	// 一見ロジックだが問題ない。
	// Query はリストを取得してしまうという DynamoDB 固有の特性をインターフェースの制約にあうようにしているだけ
	if len(manifests) < 1 {
		return dto.FindManifestOutput{}, nil
	}
	manifest := manifests[0]
	decordedManifest, err := base64.StdEncoding.DecodeString(manifest.Manifest)
	if err != nil {
		return dto.FindManifestOutput{}, err
	}

	return dto.FindManifestOutput{
		Name:     manifest.Name,
		Tag:      manifest.Tag,
		Digest:   manifest.Digest,
		Manifest: decordedManifest,
	}, nil
}

func (r ManifestRepository) SaveManifest(input dto.SaveManifestInput) error {
	base64Manifest := base64.StdEncoding.EncodeToString(input.Manifest)

	dbManifest := Manifest{
		Name:     input.Name,
		Digest:   input.Digest,
		Tag:      input.Tag,
		Manifest: base64Manifest,
	}

	item, err := attributevalue.MarshalMap(dbManifest)
	if err != nil {
		return err
	}
	_, err = r.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(r.manifestTableName),
		Item:      item,
	})
	return err
}

func (r ManifestRepository) DeleteManifest(input dto.DeleteManifestInput) error {
	if domain.IsDigest(input.Reference) {
		return r.DeleteManifestByDigest(input)
	} else {
		return r.DeleteManifestByDigest(input)
	}
}

func (r ManifestRepository) DeleteManifestByDigest(input dto.DeleteManifestInput) error {
	itemInput := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.manifestTableName),
		Key: map[string]types.AttributeValue{
			"Name": &types.AttributeValueMemberS{
				Value: input.Name,
			},
			"Digest": &types.AttributeValueMemberS{
				Value: input.Reference,
			},
		},
	}

	_, err := r.client.DeleteItem(context.TODO(), itemInput)
	return err

}

func (r ManifestRepository) DeleteManifestByTag(input dto.DeleteManifestInput) error {
	keyEx := expression.KeyAnd(
		expression.Key("Name").Equal(expression.Value(input.Name)),
		expression.Key("Tag").Equal(expression.Value(input.Reference)),
	)
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return err
	}
	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(r.manifestTableName),
		IndexName:                 aws.String("ManifestTagIndex"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}
	if err != nil {
		return err
	}
	resp, err := r.QueryItem(context.TODO(), queryInput)
	if err != nil {
		return err
	}
	var manifests []Manifest
	err = attributevalue.UnmarshalListOfMaps(resp.Items, &manifests)
	if err != nil {
		return err
	}

	if len(manifests) < 1 {
		// no such tag, but success
		return nil
	}
	manifest := manifests[0]
	return r.DeleteManifestByDigest(dto.DeleteManifestInput{
		Name:      input.Name,
		Reference: manifest.Digest,
	})
}
