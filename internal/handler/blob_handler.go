package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/a-takamin/tcr/internal/apperrors"
	"github.com/a-takamin/tcr/internal/dto"
	"github.com/a-takamin/tcr/internal/service/usecase"
	"github.com/gin-gonic/gin"
)

type BlobHandler struct {
	usecase *usecase.BlobUseCase
}

func NewBlobHandler(s *usecase.BlobUseCase) *BlobHandler {
	return &BlobHandler{
		usecase: s,
	}
}
func (h *BlobHandler) ExistsBlobHandler(c *gin.Context, name string, digest string) {
	metadata := dto.GetBlobInput{
		Name:   name,
		Digest: digest,
	}

	_, err := h.usecase.GetBlob(metadata)
	if err != nil {
		apperrors.ErrorHanlder(c, err)
		return
	}

	c.Header("Docker-Content-Digest", digest)
	c.JSON(http.StatusOK, "")
}

func (h *BlobHandler) GetBlobHandler(c *gin.Context, name string, digest string) {
	metadata := dto.GetBlobInput{
		Name:   name,
		Digest: digest,
	}

	blob, err := h.usecase.GetBlob(metadata)
	if err != nil {
		apperrors.ErrorHanlder(c, err)
		return
	}

	c.Header("Docker-Content-Digest", digest)
	c.JSON(http.StatusOK, blob)
}

func (h *BlobHandler) StartUploadBlobHandler(c *gin.Context, name string) {
	redirectUrl, err := h.usecase.StartBlobUpload(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.Header("Location", redirectUrl)
	c.JSON(http.StatusAccepted, "")
}

func (h *BlobHandler) UploadBlobHandler(c *gin.Context, name string, uuid string) {
	digest := c.Query("digest")
	ContentLength := c.Request.ContentLength
	ContentRange := c.Request.Header.Get("Content-Range")
	ContentType := c.ContentType()
	bodyStream := c.Request.Body

	isChunkedUpload, err := h.usecase.IsChunkedUpload(name, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
		return
	}
	if !isChunkedUpload {
		// Monolithic Upload
		input := dto.UploadMonolithicBlobInput{
			Name:          name,
			Uuid:          uuid,
			Digest:        digest,
			ContentLength: ContentLength,
			ContentType:   ContentType,
			Blob:          bodyStream,
		}
		err := h.usecase.UploadMonolithicBlob(input)
		if err != nil {
			// TODO: http status code
			slog.Error(err.Error())
			c.JSON(http.StatusInternalServerError, "")
			return
		}
		// TODO: http status code
		c.Header("Location", fmt.Sprintf("/v2/%s/blobs/%s", name, digest))
		c.JSON(http.StatusCreated, "")
		return
	}

	// Last Chunked Blob Upload
	input := dto.UploadChunkedBlobInput{
		Name:          name,
		Uuid:          uuid,
		Digest:        digest,
		ContentLength: ContentLength,
		ContentRange:  ContentRange,
		ContentType:   ContentType,
		Blob:          bodyStream,
		IsLast:        true,
	}

	offset, err := h.usecase.UploadLastChunkedBlob(input)

	if err != nil {
		slog.Error(err.Error())
		c.Header("Location", c.Request.URL.Path)
		c.Header("Content-Length", "0")
		c.Header("Docker-Upload-UUID", uuid)
		c.Header("Range", fmt.Sprintf("0-%d", offset))
		c.JSON(http.StatusRequestedRangeNotSatisfiable, "")
		return
	}

	c.Header("Location", fmt.Sprintf("/v2/%s/blobs/%s", name, digest))
	c.Header("Content-Length", "0")
	c.Header("Docker-Upload-Digest", digest)
	c.JSON(http.StatusCreated, "")
}

func (h *BlobHandler) UploadChunkedBlobHandler(c *gin.Context, name string, uuid string) {
	ContentLength := c.Request.ContentLength
	ContentRange := c.Request.Header.Get("Content-Range")
	if ContentRange == "" {
		ContentRange = fmt.Sprintf("0-%d", ContentLength)
	}
	ContentType := c.ContentType()
	bodyStream := c.Request.Body

	input := dto.UploadChunkedBlobInput{
		Name:          name,
		Uuid:          uuid,
		ContentLength: ContentLength,
		ContentRange:  ContentRange,
		ContentType:   ContentType,
		Blob:          bodyStream,
	}

	offset, err := h.usecase.UploadChunkedBlob(input)

	c.Header("Location", c.Request.URL.Path)
	c.Header("Content-Length", "0")
	c.Header("Docker-Upload-UUID", uuid)

	if err != nil {
		slog.Error(err.Error())
		c.Header("Range", fmt.Sprintf("0-%d", offset))
		c.JSON(http.StatusRequestedRangeNotSatisfiable, "")
	}

	c.Header("Range", fmt.Sprintf("0-%d", offset))
	c.JSON(http.StatusAccepted, "")
}

func (h *BlobHandler) DeleteBlobHandler(c *gin.Context, name string, digest string) {
	input := dto.DeleteBlobInput{
		Name:   name,
		Digest: digest,
	}

	err := h.usecase.DeleteBlob(input)
	if err != nil {
		slog.Error(err.Error())
		apperrors.ErrorHanlder(c, err)
		return
	}

	c.JSON(http.StatusAccepted, "")
}

func (h *BlobHandler) GetUploadStatusHandler(c *gin.Context, name string, uuid string) {
	offset, err := h.usecase.GetBlobUploadStatus(name, uuid)
	if err != nil {
		slog.Error(err.Error())
		// TODO: http status
		c.JSON(http.StatusNotFound, "")
		return
	}
	c.Header("Range", fmt.Sprintf("0-%d", offset))
	c.Header("Content-Length", "0")
	c.Header("Blob-Upload-Session-ID", uuid)
	c.Header("Location", fmt.Sprintf("/v2/%s/blobs/uploads/%s", name, uuid))
	c.JSON(http.StatusNoContent, "")
}
