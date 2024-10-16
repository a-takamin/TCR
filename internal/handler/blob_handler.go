package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

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
	metadata := dto.FindBlobInput{
		Name:   name,
		Digest: digest,
	}

	_, err := h.usecase.ExistsBlob(metadata)
	if err != nil {
		slog.Error(err.Error())
		switch {
		case errors.Is(err, apperrors.TCRERR_NAME_INVALID), errors.Is(err, apperrors.TCRERR_DIGEST_INVALID):
			c.JSON(http.StatusBadRequest, apperrors.NAME_INVALID.CreateResponse(""))
		case errors.Is(err, apperrors.TCRERR_NAME_NOT_FOUND), errors.Is(err, apperrors.TCRERR_BLOB_NOT_FOUND):
			c.JSON(http.StatusNotFound, apperrors.NAME_UNKNOWN.CreateResponse(""))
		case errors.Is(err, apperrors.TCRERR_PERSISTER_ERROR):
			c.JSON(http.StatusInternalServerError, "")
		case errors.Is(err, apperrors.TCRERR_LOGIC_ERROR):
			c.JSON(http.StatusInternalServerError, "")
		default:
			c.JSON(http.StatusInternalServerError, "")
		}
		return
	}

	c.Header("Docker-Content-Digest", digest)
	c.JSON(http.StatusOK, "")
}

func (h *BlobHandler) GetBlobHandler(c *gin.Context, name string, digest string) {
	metadata := dto.FindBlobInput{
		Name:   name,
		Digest: digest,
	}

	blob, err := h.usecase.GetBlob(metadata)
	if err != nil {
		slog.Error(err.Error())
		switch {
		case errors.Is(err, apperrors.TCRERR_NAME_INVALID), errors.Is(err, apperrors.TCRERR_DIGEST_INVALID):
			c.JSON(http.StatusBadRequest, apperrors.NAME_INVALID.CreateResponse(""))
		case errors.Is(err, apperrors.TCRERR_NAME_NOT_FOUND), errors.Is(err, apperrors.TCRERR_BLOB_NOT_FOUND):
			c.JSON(http.StatusNotFound, apperrors.NAME_UNKNOWN.CreateResponse(""))
		case errors.Is(err, apperrors.TCRERR_PERSISTER_ERROR):
			c.JSON(http.StatusInternalServerError, "")
		case errors.Is(err, apperrors.TCRERR_LOGIC_ERROR):
			c.JSON(http.StatusInternalServerError, "")
		default:
			c.JSON(http.StatusInternalServerError, "")
		}
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

// Put はモノリスとチャンクのラストとの 2 通りがある
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
	ContentLengthHeaderVal := c.Request.Header.Get("Content-Length")
	var ContentLength int64
	if ContentLengthHeaderVal == "" {
		ContentLength = 0
	} else {
		var err error
		ContentLength, err = strconv.ParseInt(ContentLengthHeaderVal, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
	}
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
		return
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
	offset, err := h.usecase.GetBlobUploadOffset(name, uuid)
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
