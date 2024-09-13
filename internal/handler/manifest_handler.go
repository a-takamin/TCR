package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/a-takamin/tcr/internal/apperrors"
	"github.com/a-takamin/tcr/internal/model"
	"github.com/a-takamin/tcr/internal/service"
	"github.com/a-takamin/tcr/internal/service/utils"
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
	name := c.Param("name")
	reference := c.Param("reference")

	metadata := model.ManifestMetadata{
		Name:      name,
		Reference: reference,
	}

	manifest, err := h.service.GetManifest(metadata)
	if err != nil {
		if errors.Is(err, apperrors.ErrManifestNotFound) {
			c.JSON(http.StatusNotFound, err)
			return
		}
		c.JSON(http.StatusBadRequest, err)
		return
	}

	digest, err := utils.CalcManifestDigest(manifest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.Header("Docker-Content-Digest", digest)
	c.JSON(http.StatusOK, manifest)
}

func (h *ManifestHandler) PutManifestHandler(c *gin.Context) {
	name := c.Param("name")
	reference := c.Param("reference")

	metadata := model.ManifestMetadata{
		Name:      name,
		Reference: reference,
	}

	contentType := c.Request.Header.Get("Content-Type")

	var manifest model.Manifest

	if err := c.ShouldBindJSON(&manifest); err != nil {
		err := fmt.Errorf("manifest is invalid: %w", err)
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if contentType != manifest.MediaType {
		err := errors.New("Content-Type is invalid")
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// TODO: manifest が指す Blob があるかどうか MUST で確認する。なければ 404 を返す。

	err := h.service.PutManifest(metadata, manifest)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.Redirect(http.StatusCreated, c.Request.Host+c.Request.URL.Path)
}

func (h *ManifestHandler) DeleteManifestHandler(c *gin.Context) {
	name := c.Param("name")
	reference := c.Param("reference")

	metadata := model.ManifestMetadata{
		Name:      name,
		Reference: reference,
	}

	err := h.service.DeleteManifest(metadata)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusAccepted, "")
}
