package main

import (
	"github.com/a-takamin/tcr/handler"
	"github.com/a-takamin/tcr/repository"
	"github.com/a-takamin/tcr/service"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	repo := repository.NewDynamoDBManifestRepository()
	s := service.NewManifestService(repo)

	h := handler.NewManifestHandler(s)

	r.GET("/v2/:name/manifests/:reference", h.GetManifestHandler)

	r.Run(":8080")
}
