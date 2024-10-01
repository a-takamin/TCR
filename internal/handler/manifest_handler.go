package handler

import (
	"fmt"
	"io"
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

func (h *ManifestHandler) ExistsManifestHandler(c *gin.Context, name string, reference string) {
	metadata := model.ManifestMetadata{
		Name:      name,
		Reference: reference,
	}

	resp, err := h.usecase.ExistsManifest(metadata)
	if err != nil {
		apperrors.ErrorHanlder(c, err)
		return
	}

	c.Header("Docker-Content-Digest", resp.Digest)
	c.JSON(http.StatusOK, "")
}

func (h *ManifestHandler) GetManifestHandler(c *gin.Context, name string, reference string) {
	metadata := model.ManifestMetadata{
		Name:      name,
		Reference: reference,
	}

	resp, err := h.usecase.GetManifest(metadata)
	if err != nil {
		apperrors.ErrorHanlder(c, err)
		return
	}

	c.Header("Docker-Content-Digest", resp.Digest)
	c.JSON(http.StatusOK, resp.Manifest)
}

func (h *ManifestHandler) PutManifestHandler(c *gin.Context, name string, reference string) {
	metadata := model.ManifestMetadata{
		Name:        name,
		Reference:   reference,
		ContentType: c.Request.Header.Get("Content-Type"),
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		err := fmt.Errorf("could not read manifest: %w", err)
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// TODO: manifest が指す Blob があるかどうか MUST で確認する。なければ 404 を返す。

	err = h.usecase.PutManifest(metadata, body)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.Redirect(http.StatusCreated, c.Request.Host+c.Request.URL.Path)
}

func (h *ManifestHandler) DeleteManifestHandler(c *gin.Context, name string, reference string) {
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
