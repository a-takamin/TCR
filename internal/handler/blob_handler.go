package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/a-takamin/tcr/apperrors"
	"github.com/a-takamin/tcr/internal/model"
	"github.com/a-takamin/tcr/internal/service"
	"github.com/a-takamin/tcr/internal/service/utils"
	"github.com/gin-gonic/gin"
)

type BlobHandler struct {
	service *service.BlobService
}

func NewBlobHandler(s *service.BlobService) *BlobHandler {
	return &BlobHandler{
		service: s,
	}
}

func (h *BlobHandler) GetBlobHandler(c *gin.Context) {
	name := c.Param("name")
	digest := c.Param("digest")

	metadata := model.BlobMetadata{
		Name:   name,
		Digest: digest,
	}

	blob, err := h.service.GetBlob(metadata)
	if err != nil {
		if errors.Is(err, apperrors.ErrBlobNotFound) {
			c.JSON(http.StatusNotFound, err)
			return
		}
		c.JSON(http.StatusBadRequest, err)
		return
	}

	blobDigest, err := utils.CalcBlobDigest(blob)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.Header("Docker-Content-Digest", blobDigest)
	c.JSON(http.StatusOK, blob)
}

func (h *BlobHandler) StartUploadBlobHandler(c *gin.Context) {
	name := c.Param("name")
	redirectUrl, err := h.service.StartBlobUpload(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.Header("Location", redirectUrl)
	c.JSON(http.StatusAccepted, "")
}

func (h *BlobHandler) UploadBlobHandler(c *gin.Context) {
	name := c.Param("name")
	uuid := c.Param("uuid")
	digest := c.Query("digest")
	bodyStream := c.Request.Body

	metadata := model.BlobUploadMetadata{
		Name:          name,
		Uuid:          uuid,
		Digest:        digest,
		ContentLength: c.Request.ContentLength,
		ContentRange:  c.Request.Header.Get("Content-Range"),
		ContentType:   c.ContentType(),
		IsChunkUpload: false,
	}

	var err error

	if metadata.ContentRange == "" {
		err = h.service.UploadBlob(metadata, bodyStream)
	} else {
		// err = h.service.CompleteUploadBlob()
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, "")
	}
	c.JSON(http.StatusCreated, "")
}

func (h *BlobHandler) UploadChunkedBlobHandler(c *gin.Context) {
	name := c.Param("name")
	uuid := c.Param("uuid")
	body := c.Request.Body

	metadata := model.BlobUploadMetadata{
		Name:          name,
		Uuid:          uuid,
		ContentLength: c.Request.ContentLength,
		ContentRange:  c.Request.Header.Get("Content-Range"),
		ContentType:   c.ContentType(),
		IsChunkUpload: true,
	}

	c.Header("Location", c.Request.URL.Path)
	c.Header("Content-Length", "0")
	c.Header("Docker-Upload-UUID", uuid)

	offset, err := h.service.UploadChunkedBlob(metadata, body)
	if err != nil {
		c.Header("Range", fmt.Sprintf("0-%d", offset))
		c.JSON(http.StatusRequestedRangeNotSatisfiable, "")
	}

	c.Header("Range", fmt.Sprintf("bytes=0-%d", offset))
	c.JSON(http.StatusAccepted, "")
}
