package main

import (
	"log"
	"net/http"

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

	facade := handler.NewFacadeHandler(mh, bh)

	r.GET("/v2", func(c *gin.Context) { // end-1
		c.JSON(http.StatusOK, "")
	})

	r.HEAD("/v2/*remain", facade.HandleHEAD)     // end-2, end-3
	r.GET("/v2/*remain", facade.HandleGET)       // end-2, end-3, end-8a
	r.POST("/v2/*remain", facade.HandlePOST)     // end-4a, 4b
	r.PUT("/v2/*remain", facade.HandlePUT)       // end-6,
	r.PATCH("/v2/*remain", facade.HandlePATCH)   // end-5
	r.DELETE("/v2/*remain", facade.HandleDELETE) // end-9, end-10

	r.Run(":8080")

}
