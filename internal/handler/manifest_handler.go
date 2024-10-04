package handler

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
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
		slog.Error(err.Error())
		switch {
		case errors.Is(err, apperrors.TCRERR_NAME_INVALID):
			c.JSON(http.StatusBadRequest, apperrors.NAME_INVALID.CreateResponse(""))
		case errors.Is(err, apperrors.TCRERR_NAME_NOT_FOUND):
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
		err := fmt.Errorf("could not read manifest body: %w", err)
		slog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, "could not read manifest body")
		return
	}

	err = h.usecase.PutManifest(metadata, body)
	if err != nil {
		slog.Error(err.Error())
		switch {
		case errors.Is(err, apperrors.TCRERR_NAME_INVALID):
			c.JSON(http.StatusBadRequest, apperrors.NAME_INVALID.CreateResponse(""))
		case errors.Is(err, apperrors.TCRERR_MANIFEST_INVALID):
			c.JSON(http.StatusBadRequest, apperrors.MANIFEST_INVALID.CreateResponse(""))
		case errors.Is(err, apperrors.TCRERR_NAME_NOT_FOUND):
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
	c.Redirect(http.StatusCreated, c.Request.Host+c.Request.URL.Path)
}

func (h *ManifestHandler) DeleteManifestHandler(c *gin.Context, name string, reference string) {
	metadata := model.ManifestMetadata{
		Name:      name,
		Reference: reference,
	}

	err := h.usecase.DeleteManifest(metadata)
	if err != nil {
		slog.Error(err.Error())
		switch {
		case errors.Is(err, apperrors.TCRERR_NAME_INVALID):
			c.JSON(http.StatusBadRequest, apperrors.NAME_INVALID.CreateResponse(""))
		case errors.Is(err, apperrors.TCRERR_NAME_NOT_FOUND):
			c.JSON(http.StatusNotFound, apperrors.NAME_UNKNOWN.CreateResponse(""))
		case errors.Is(err, apperrors.TCRERR_MANIFEST_NOT_FOUND):
			c.JSON(http.StatusNotFound, apperrors.MANIFEST_UNKNOWN.CreateResponse(""))
		case errors.Is(err, apperrors.TCRERR_PERSISTER_ERROR):
			c.JSON(http.StatusInternalServerError, "")
		default:
			c.JSON(http.StatusInternalServerError, "")
		}
		return
	}
	c.JSON(http.StatusAccepted, "")
}
