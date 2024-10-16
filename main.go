package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/a-takamin/tcr/internal/client"
	"github.com/a-takamin/tcr/internal/handler"
	"github.com/a-takamin/tcr/internal/repository"
	"github.com/a-takamin/tcr/internal/service/usecase"
	"github.com/gin-gonic/gin"
)

func init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, nil)))
}

func main() {
	r := gin.New()
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/health"},
	}))
	r.Use(handler.LogMiddleWare())
	r.Use(gin.Recovery())
	isLocal := true

	env := os.Getenv("IS_LOCAL")
	if env != "" {
		isLocal = false
	}
	blobStorageName := os.Getenv("BLOB_STORAGE_NAME")
	if blobStorageName == "" {
		blobStorageName = "tcr-blob-local"
	}
	manifestTableName := os.Getenv("MANIFEST_TABLE_NAME")
	if manifestTableName == "" {
		manifestTableName = "tcr-manifest-local"
	}
	repositoryTableName := os.Getenv("REPOSITORY_TABLE_NAME")
	if repositoryTableName == "" {
		repositoryTableName = "tcr-repository-local"
	}
	blobUploadProgressTableName := os.Getenv("BLOB_UPLOAD_PROGRESS_TABLE_NAME")
	if blobUploadProgressTableName == "" {
		blobUploadProgressTableName = "tcr-blob-upload-progress-local"
	}

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

	mRepo := repository.NewManifestRepository(dynamodbClient, manifestTableName)
	bRepo := repository.NewBlobRepository(s3Client, blobStorageName, dynamodbClient, blobUploadProgressTableName)
	rRepo := repository.NewRepositoryRepository(dynamodbClient, repositoryTableName)
	pRepo := repository.NewBlobUploadProgressRepository(dynamodbClient, blobUploadProgressTableName)

	mu := usecase.NewManifestUseCase(mRepo, rRepo)
	bu := usecase.NewBlobUseCase(bRepo, pRepo, rRepo)

	mh := handler.NewManifestHandler(mu)
	bh := handler.NewBlobHandler(bu)

	facade := handler.NewFacadeHandler(mh, bh)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})
	r.HEAD("/v2/*remain", facade.HandleHEAD)     // end-2, end-3
	r.GET("/v2/*remain", facade.HandleGET)       // end-2, end-3, end-8a
	r.POST("/v2/*remain", facade.HandlePOST)     // end-4a, 4b
	r.PUT("/v2/*remain", facade.HandlePUT)       // end-6,
	r.PATCH("/v2/*remain", facade.HandlePATCH)   // end-5
	r.DELETE("/v2/*remain", facade.HandleDELETE) // end-9, end-10

	r.Run(":8080")

}
