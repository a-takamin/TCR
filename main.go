package main

import (
	"log"

	"github.com/a-takamin/tcr/internal/client"
	"github.com/a-takamin/tcr/internal/handler"
	"github.com/a-takamin/tcr/internal/repository"
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

	blobDomain := domain.NewBlobDomain(bRepo)
	mu := usecase.NewManifestUseCase(mRepo)
	bu := usecase.NewBlobUseCase(blobDomain, bRepo)

	mh := handler.NewManifestHandler(mu)
	bh := handler.NewBlobHandler(bu)

	r.HEAD("/v2/:name/blobs/:digest", bh.ExistsBlobHandler)               // end-2
	r.GET("/v2/:name/blobs/:digest", bh.GetBlobHandler)                   // end-2
	r.HEAD("/v2/:name/manifests/:reference", mh.ExistsManifestHandler)    // end-3
	r.GET("/v2/:name/manifests/:reference", mh.GetManifestHandler)        // end-3
	r.POST("/v2/:name/blobs/uploads", bh.StartUploadBlobHandler)          // end-4a, 4b
	r.PATCH("/v2/:name/blobs/uploads/:uuid", bh.UploadChunkedBlobHandler) // end-5
	r.PUT("/v2/:name/manifests/:reference", mh.PutManifestHandler)        // end-6
	r.GET("/v2/:name/tags/list", mh.GetTagsHandler)                       // end-8a
	r.DELETE("/v2/:name/manifests/:reference", mh.DeleteManifestHandler)  // end-9
	r.DELETE("/v2/:name/blobs/:reference", bh.DeleteBlobHandler)          // end-10

	r.PUT("/v2/:name/blobs/uploads/:uuid", bh.UploadBlobHandler)

	r.Run(":8080")

}
