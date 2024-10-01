package handler

import (
	"net/http"

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
		// TODO
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, tags)
}
