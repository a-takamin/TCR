package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/a-takamin/tcr/internal/apperrors"
	"github.com/a-takamin/tcr/internal/model"
	"github.com/a-takamin/tcr/internal/service/usecase"
	"github.com/gin-gonic/gin"
)

type ManifestHandler struct {
	usecase *usecase.ManifestUseCase
}

func NewManifestHandler(u *usecase.ManifestUseCase) *ManifestHandler {
	return &ManifestHandler{
		usecase: u,
	}
}

func (h *ManifestHandler) GetManifestHandler(c *gin.Context) {
	name := c.Param("name")
	reference := c.Param("reference")

	metadata := model.ManifestMetadata{
		Name:      name,
		Reference: reference,
	}

	resp, err := h.usecase.GetManifest(metadata)
	if err != nil {
		// TODO: エラーハンドリング
		if errors.Is(err, apperrors.ErrManifestNotFound) {
			c.JSON(http.StatusNotFound, err)
			return
		}
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.Header("Docker-Content-Digest", resp.Digest)
	c.JSON(http.StatusOK, resp.Manifest)
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

	err := h.usecase.PutManifest(metadata, manifest)
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

	err := h.usecase.DeleteManifest(metadata)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusAccepted, "")
}

func (h *ManifestHandler) GetTagsHandler(c *gin.Context) {
	name := c.Param("name")
	tags, err := h.usecase.GetTags(name)
	if err != nil {
		// TODO
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, tags)
}
