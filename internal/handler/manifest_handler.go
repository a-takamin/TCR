package handler

import (
	"github.com/a-takamin/tcr/internal/service"
	"github.com/gin-gonic/gin"
)

type ManifestHandler struct {
	service *service.ManifestService
}

func NewManifestHandler(s *service.ManifestService) *ManifestHandler {
	return &ManifestHandler{
		service: s,
	}
}

func (h *ManifestHandler) GetManifestHandler(c *gin.Context) {
	// metadata := models.ManifestMetadata{}
	// manifest, err := h.service.GetManifest(metadata)
	// _ = manifest
	// _ = err
}
