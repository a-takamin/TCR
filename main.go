package main

import (
	"log"

	"github.com/a-takamin/tcr/internal/client"
	"github.com/a-takamin/tcr/internal/handler"
	"github.com/a-takamin/tcr/internal/repository"
	"github.com/a-takamin/tcr/internal/service"
	"github.com/a-takamin/tcr/internal/service/domain"
	"github.com/a-takamin/tcr/internal/service/usecase"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	isLocal := true

	dynamodbClient, err := client.NewDynamoDbClient(isLocal)
	if err != nil {
		log.Fatal(err)
		return
	}

	s3Client, err := client.NewS3Client(isLocal)
	if err != nil {
		log.Fatal(err)
		return
	}

	mRepo := repository.NewManifestRepository(dynamodbClient, "dynamodb-local-table")
	bRepo := repository.NewBlobRepository(s3Client, "blob-local", dynamodbClient, "blob-upload-progress")
	ms := service.NewManifestService(mRepo)
	bs := domain.NewBlobDomain(bRepo)
	// mu := usecase.NewManifestUseCase(ms)
	bu := usecase.NewBlobUseCase(bs, bRepo)

	mh := handler.NewManifestHandler(ms)
	bh := handler.NewBlobHandler(bu)

	r.GET("/v2/:name/manifests/:reference", mh.GetManifestHandler)
	r.PUT("/v2/:name/manifests/:reference", mh.PutManifestHandler)
	r.DELETE("/v2/:name/manifests/:reference", mh.DeleteManifestHandler)

	r.GET("/v2/:name/blobs/:digest", bh.GetBlobHandler)
	r.POST("/v2/:name/blobs/uploads", bh.StartUploadBlobHandler)
	r.PUT("/v2/:name/blobs/uploads/:uuid", bh.UploadBlobHandler)
	r.PATCH("/v2/:name/blobs/uploads/:uuid", bh.UploadChunkedBlobHandler)

	r.GET("/v2/:name/tags/list", bh.GetTagsHandler)

	r.Run(":8080")

}
