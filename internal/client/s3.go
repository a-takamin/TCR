package client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewS3Client(isLocal bool) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-1"))
	if err != nil {
		return nil, err
	}
	if isLocal {
		return s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String("http://localhost:9000")
			o.UsePathStyle = true // local のときだけ minio を使うために必要
			o.EndpointOptions.DisableHTTPS = true
		}), nil
	}
	return s3.NewFromConfig(cfg), nil
}
