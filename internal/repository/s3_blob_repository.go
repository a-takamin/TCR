package repository

import (
	"context"
	"io"

	"github.com/a-takamin/tcr/internal/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type BlobRepository struct {
	client     *s3.Client
	bucketName string
}

type Blob struct {
	Digest string
	Tag    string
	Blob   string
}

func NewBlobRepository(client *s3.Client, tableName string) *BlobRepository {
	return &BlobRepository{
		client:     client,
		bucketName: tableName,
	}
}

func (r BlobRepository) GetBlob(metadata model.BlobMetadata) (model.Blob, error) {
	return model.Blob{}, nil
}

func (r BlobRepository) UploadBlob(metadata model.BlobUploadMetadata, blob io.Reader) error {
	_, err := r.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(metadata.Key),
	})
	return err
}

// func (r BlobRepository) getItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
// 	return r.client.GetItem(ctx, params, optFns...)
// }

// func (r BlobRepository) getItemByTag(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
// 	return r.client.Query(ctx, params, optFns...)
// }

// func (r BlobRepository) createBlobGetInput(digest string) *dynamodb.GetItemInput {
// 	return &dynamodb.GetItemInput{
// 		TableName: aws.String(r.tableName),
// 		Key: map[string]types.AttributeValue{
// 			"Digest": &types.AttributeValueMemberS{
// 				Value: digest,
// 			},
// 		},
// 	}
// }

// func (r BlobRepository) createBlobGetResponse(manifest Blob) (model.Blob, error) {
// 	decordedBlob, err := base64.StdEncoding.DecodeString(manifest.Blob)
// 	if err != nil {
// 		return model.Blob{}, err
// 	}

// 	var modelBlob model.Blob
// 	err = json.Unmarshal(decordedBlob, &modelBlob)
// 	if err != nil {
// 		return model.Blob{}, err
// 	}

// 	return modelBlob, nil
// }

// func (r BlobRepository) GetBlob(metadata model.BlobMetadata) (model.Blob, error) {
// 	if !utils.IsDigest(metadata.Reference) {
// 		return r.GetBlobByTag(metadata)
// 	}

// 	input := r.createBlobGetInput(metadata.Reference)

// 	resp, err := r.getItem(context.TODO(), input)

// 	if err != nil {
// 		return model.Blob{}, err
// 	}

// 	if len(resp.Item) == 0 {
// 		return model.Blob{}, apperrors.ErrBlobNotFound
// 	}

// 	var manifest Blob
// 	err = attributevalue.UnmarshalMap(resp.Item, &manifest)
// 	if err != nil {
// 		return model.Blob{}, err
// 	}

// 	return r.createBlobGetResponse(manifest)
// }

// func (r BlobRepository) createBlobGetInputByTag(tag string) (*dynamodb.QueryInput, error) {
// 	keyEx := expression.Key("Tag").Equal(expression.Value(tag))
// 	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
// 	if err != nil {
// 		return &dynamodb.QueryInput{}, err
// 	}
// 	return &dynamodb.QueryInput{
// 		TableName:                 aws.String(r.tableName),
// 		ExpressionAttributeNames:  expr.Names(),
// 		ExpressionAttributeValues: expr.Values(),
// 		KeyConditionExpression:    expr.KeyCondition(),
// 		IndexName:                 aws.String("BlobDigestIndex"),
// 	}, nil
// }

// func (r BlobRepository) GetBlobByTag(metadata model.BlobMetadata) (model.Blob, error) {
// 	input, err := r.createBlobGetInputByTag(metadata.Reference)
// 	if err != nil {
// 		return model.Blob{}, err
// 	}
// 	resp, err := r.getItemByTag(context.TODO(), input)

// 	if err != nil {
// 		return model.Blob{}, err
// 	}
// 	var manifests []Blob
// 	err = attributevalue.UnmarshalListOfMaps(resp.Items, &manifests)
// 	if err != nil {
// 		return model.Blob{}, err
// 	}

// 	if len(manifests) < 1 {
// 		return model.Blob{}, apperrors.ErrBlobNotFound
// 	}
// 	manifest := manifests[0]
// 	return r.createBlobGetResponse(manifest)
// }

// func (r BlobRepository) PutBlob(metadata model.BlobMetadata, content model.Blob) error {
// 	byteBlob, err := json.Marshal(content)
// 	if err != nil {
// 		return err
// 	}
// 	encodedBlob := base64.StdEncoding.EncodeToString(byteBlob)

// 	var manifest Blob
// 	if utils.IsDigest(metadata.Reference) {
// 		manifest = Blob{
// 			Digest: metadata.Reference,
// 			Blob:   encodedBlob,
// 		}
// 	} else {
// 		digest, err := utils.CalcBlobDigest(content)
// 		if err != nil {
// 			return err
// 		}
// 		manifest = Blob{
// 			Digest: digest,
// 			Tag:    metadata.Reference,
// 			Blob:   encodedBlob,
// 		}
// 	}

// 	item, err := attributevalue.MarshalMap(manifest)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = r.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
// 		TableName: aws.String(r.tableName),
// 		Item:      item,
// 	})

// 	return err
// }

// func (r BlobRepository) DeleteBlob(metadata model.BlobMetadata) error {
// 	if !utils.IsDigest(metadata.Reference) {
// 		return r.DeleteBlobByTag(metadata)
// 	}

// 	input := &dynamodb.DeleteItemInput{
// 		TableName: aws.String(r.tableName),
// 		Key: map[string]types.AttributeValue{
// 			"Digest": &types.AttributeValueMemberS{
// 				Value: metadata.Reference,
// 			},
// 		},
// 	}

// 	_, err := r.client.DeleteItem(context.TODO(), input)
// 	return err

// }

// func (r BlobRepository) DeleteBlobByTag(metadata model.BlobMetadata) error {
// 	input, err := r.createBlobGetInputByTag(metadata.Reference)
// 	if err != nil {
// 		return err
// 	}
// 	resp, err := r.getItemByTag(context.TODO(), input)
// 	if err != nil {
// 		return err
// 	}
// 	var manifests []Blob
// 	err = attributevalue.UnmarshalListOfMaps(resp.Items, &manifests)
// 	if err != nil {
// 		return err
// 	}

// 	if len(manifests) < 1 {
// 		// TODO: make an error code
// 		return errors.New("no manifest exists")
// 	}
// 	manifest := manifests[0]
// 	return r.DeleteBlob(model.BlobMetadata{
// 		Name:      metadata.Name,
// 		Reference: manifest.Digest,
// 	})
// }
