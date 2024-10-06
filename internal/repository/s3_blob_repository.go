package repository

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/a-takamin/tcr/internal/dto"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Type "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type BlobRepository struct {
	client                      *s3.Client
	dClient                     *dynamodb.Client
	bucketName                  string
	blobUploadProgressTableName string
}

type Blob struct {
	Digest string
	Tag    string
	Blob   string
}

func NewBlobRepository(client *s3.Client, bucketName string, dynamodbClient *dynamodb.Client, blobUploadProgressTName string) *BlobRepository {
	return &BlobRepository{
		client:                      client,
		bucketName:                  bucketName,
		dClient:                     dynamodbClient,
		blobUploadProgressTableName: blobUploadProgressTName,
	}
}

// Refactor
func (r BlobRepository) ExistsBlob(input dto.ExistsBlobInput) (bool, error) {
	_, err := r.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(input.Name + "/" + input.Digest),
	})
	// GetObject はオブジェクトがないときにエラーを返す
	if err != nil {
		var noSuchKeyErr *s3Type.NoSuchKey
		if errors.As(err, &noSuchKeyErr) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r BlobRepository) FindBlob(input dto.FindBlobInput) (dto.FindBlobOutput, error) {
	resp, err := r.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(input.Name + "/" + input.Digest),
	})
	if err != nil {
		return dto.FindBlobOutput{}, err
	}
	blob, err := io.ReadAll(resp.Body)
	if err != nil {
		return dto.FindBlobOutput{}, err
	}

	return dto.FindBlobOutput{
		Blob: blob,
	}, nil
}

func (r BlobRepository) FindChunkedBlob(input dto.FindChunkedBlobInput) (dto.FindBlobOutput, error) {
	resp, err := r.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(fmt.Sprintf("/%s/chunk/%s/%d", input.Name, input.Uuid, input.ChunkSeqNo)),
	})
	if err != nil {
		return dto.FindBlobOutput{}, err
	}
	blob, err := io.ReadAll(resp.Body)
	if err != nil {
		return dto.FindBlobOutput{}, err
	}

	return dto.FindBlobOutput{
		Blob: blob,
	}, nil
}

func (r BlobRepository) SaveBlob(input dto.SaveBlobInput) error {
	// TODO: なぜか分からないが「"failed to seek body to start, request stream is not seekable”」が発生するので、
	// 一度Blobを読み込んで再度Reader型にしている
	b, err := io.ReadAll(input.Blob)
	if err != nil {
		return err
	}
	_, err = r.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(input.Name + "/" + input.Digest),
		Body:   bytes.NewReader(b),
	})
	return err
}

func (r BlobRepository) SaveChunkedBlob(input dto.SaveChunkedBlobInput) error {
	// TODO: なぜか分からないが「"failed to seek body to start, request stream is not seekable”」が発生するので、
	// 一度Blobを読み込んで再度Reader型にしている
	b, err := io.ReadAll(input.Blob)
	if err != nil {
		return err
	}
	_, err = r.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(fmt.Sprintf("/%s/chunk/%s/%d", input.Name, input.Uuid, input.ChunkSeqNo)),
		Body:   bytes.NewReader(b),
	})
	return err
}

func (r BlobRepository) DeleteBlob(input dto.DeleteBlobInput) error {
	_, err := r.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(input.Name + "/" + input.Digest),
	})
	return err
}
