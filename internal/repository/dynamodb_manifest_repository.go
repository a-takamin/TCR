package repository

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/a-takamin/tcr/apperrors"
	"github.com/a-takamin/tcr/internal/model"
	"github.com/a-takamin/tcr/internal/service/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type ManifestRepository struct {
	client    *dynamodb.Client
	tableName string
}

type Manifest struct {
	Digest   string `dynamodbav:Digest`
	Tag      string `dynamodbav:Tag`
	Manifest string `dynamodbav:Manifest`
}

func NewManifestRepository(client *dynamodb.Client, tableName string) *ManifestRepository {
	return &ManifestRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r ManifestRepository) getItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	return r.client.GetItem(ctx, params, optFns...)
}

func (r ManifestRepository) getItemByTag(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	return r.client.Query(ctx, params, optFns...)
}

func (r ManifestRepository) createManifestGetInput(digest string) *dynamodb.GetItemInput {
	return &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"Digest": &types.AttributeValueMemberS{
				Value: digest,
			},
		},
	}
}

func (r ManifestRepository) createManifestGetResponse(manifest Manifest) (model.Manifest, error) {
	decordedManifest, err := base64.StdEncoding.DecodeString(manifest.Manifest)
	if err != nil {
		return model.Manifest{}, err
	}

	var modelManifest model.Manifest
	err = json.Unmarshal(decordedManifest, &modelManifest)
	if err != nil {
		return model.Manifest{}, err
	}

	return modelManifest, nil
}

func (r ManifestRepository) GetManifest(metadata model.ManifestMetadata) (model.Manifest, error) {
	if !utils.IsDigest(metadata.Reference) {
		return r.GetManifestByTag(metadata)
	}

	input := r.createManifestGetInput(metadata.Reference)

	resp, err := r.getItem(context.TODO(), input)

	if err != nil {
		return model.Manifest{}, err
	}

	if len(resp.Item) == 0 {
		return model.Manifest{}, apperrors.ErrManifestNotFound
	}

	var manifest Manifest
	err = attributevalue.UnmarshalMap(resp.Item, &manifest)
	if err != nil {
		return model.Manifest{}, err
	}

	return r.createManifestGetResponse(manifest)
}

func (r ManifestRepository) createManifestGetInputByTag(tag string) (*dynamodb.QueryInput, error) {
	keyEx := expression.Key("Tag").Equal(expression.Value(tag))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return &dynamodb.QueryInput{}, err
	}
	return &dynamodb.QueryInput{
		TableName:                 aws.String(r.tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		IndexName:                 aws.String("ManifestTagIndex"),
	}, nil
}

func (r ManifestRepository) GetManifestByTag(metadata model.ManifestMetadata) (model.Manifest, error) {
	input, err := r.createManifestGetInputByTag(metadata.Reference)
	if err != nil {
		return model.Manifest{}, err
	}
	resp, err := r.getItemByTag(context.TODO(), input)

	if err != nil {
		return model.Manifest{}, err
	}
	var manifests []Manifest
	err = attributevalue.UnmarshalListOfMaps(resp.Items, &manifests)
	if err != nil {
		return model.Manifest{}, err
	}

	if len(manifests) < 1 {
		return model.Manifest{}, apperrors.ErrManifestNotFound
	}
	manifest := manifests[0]
	return r.createManifestGetResponse(manifest)
}

func (r ManifestRepository) PutManifest(metadata model.ManifestMetadata, content model.Manifest) error {
	byteManifest, err := json.Marshal(content)
	if err != nil {
		return err
	}
	encodedManifest := base64.StdEncoding.EncodeToString(byteManifest)

	var manifest Manifest
	if utils.IsDigest(metadata.Reference) {
		manifest = Manifest{
			Digest:   metadata.Reference,
			Manifest: encodedManifest,
		}
	} else {
		digest, err := utils.CalcManifestDigest(content)
		if err != nil {
			return err
		}
		manifest = Manifest{
			Digest:   digest,
			Tag:      metadata.Reference,
			Manifest: encodedManifest,
		}
	}

	item, err := attributevalue.MarshalMap(manifest)
	if err != nil {
		return err
	}
	_, err = r.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})

	return err
}

func (r ManifestRepository) DeleteManifest(metadata model.ManifestMetadata) error {
	if !utils.IsDigest(metadata.Reference) {
		return r.DeleteManifestByTag(metadata)
	}

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"Digest": &types.AttributeValueMemberS{
				Value: metadata.Reference,
			},
		},
	}

	_, err := r.client.DeleteItem(context.TODO(), input)
	return err

}

func (r ManifestRepository) DeleteManifestByTag(metadata model.ManifestMetadata) error {
	input, err := r.createManifestGetInputByTag(metadata.Reference)
	if err != nil {
		return err
	}
	resp, err := r.getItemByTag(context.TODO(), input)
	if err != nil {
		return err
	}
	var manifests []Manifest
	err = attributevalue.UnmarshalListOfMaps(resp.Items, &manifests)
	if err != nil {
		return err
	}

	if len(manifests) < 1 {
		// TODO: make an error code
		return errors.New("no manifest exists")
	}
	manifest := manifests[0]
	return r.DeleteManifest(model.ManifestMetadata{
		Name:      metadata.Name,
		Reference: manifest.Digest,
	})
}
