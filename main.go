package main

import (
	"context"
	"log"

	"github.com/a-takamin/tcr/internal/handler"
	"github.com/a-takamin/tcr/internal/repository"
	"github.com/a-takamin/tcr/internal/service"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// TODO: DynamoDB Client の作成を別パッケージにする
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal()
	}
	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String("http://localhost:8000")
	})

	repo := repository.NewManifestRepository(client, "dynamodb-local-table")
	s := service.NewManifestService(repo)

	h := handler.NewManifestHandler(s)

	r.GET("/v2/:name/manifests/:reference", h.GetManifestHandler)
	r.PUT("/v2/:name/manifests/:reference", h.PutManifestHandler)
	r.DELETE("/v2/:name/manifests/:reference", h.DeleteManifestHandler)

	r.Run(":8080")

}
