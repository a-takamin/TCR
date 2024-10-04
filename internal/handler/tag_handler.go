package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/a-takamin/tcr/internal/apperrors"
	"github.com/a-takamin/tcr/internal/service/usecase"
	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	usecase *usecase.TagUseCase
}

func NewTagHandler(u *usecase.TagUseCase) *TagHandler {
	return &TagHandler{
		usecase: u,
	}
}

func (h TagHandler) GetTagsHandler(c *gin.Context, name string) {
	tags, err := h.usecase.GetTags(name)
	if err != nil {
		slog.Error(err.Error())
		switch {
		case errors.Is(err, apperrors.TCRERR_PERSISTER_ERROR):
			c.JSON(http.StatusInternalServerError, "")
		case errors.Is(err, apperrors.TCRERR_NAME_NOT_FOUND):
			c.JSON(http.StatusNotFound, apperrors.NAME_UNKNOWN.CreateResponse(""))
		default:
			c.JSON(http.StatusInternalServerError, "")
		}
		return
	}
	c.JSON(http.StatusOK, tags)
}
