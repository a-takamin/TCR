package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Gin ではパスの変数にスラッシュを使えないために設けられたハンドラー
type FacadeHandler struct {
	blobHandler     BlobHandler
	manifestHandler ManifestHandler
}

func NewFacadeHandler() *FacadeHandler {
	return &FacadeHandler{}
}

// 以下の API
//
// "/v2/:name/blobs/:digest"
//
// "/v2/:name/manifests/:reference"
func (h FacadeHandler) HandleHEAD(c *gin.Context) {
	remainPath := c.Param("remain")
	name, afterNamePath, err := pickUpName(remainPath, 2)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusNotFound, "")
		return
	}

	afterNameParts := strings.Split(afterNamePath, "/")
	if len(afterNameParts) != 3 {
		slog.Error("path is invalid")
		c.JSON(http.StatusNotFound, "")
		return
	}

	category := afterNameParts[1]
	reference := afterNameParts[2]

	switch category {
	case "blobs":
		h.blobHandler.ExistsBlobHandler(c, name, reference)
	case "manifests":
		h.manifestHandler.ExistsManifestHandler(c, name, reference)
	default:
		slog.Error("path is invalid: " + remainPath)
		c.JSON(http.StatusNotFound, "")
	}
}

// 以下の API
//
// "/v2/:name/blobs/:digest"
//
// "/v2/:name/manifests/:reference"
//
// "/v2/:name/tags/list"
func (h FacadeHandler) HandleGET(c *gin.Context) {
	remainPath := c.Param("remain")
	name, afterNamePath, err := pickUpName(remainPath, 2)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusNotFound, "")
		return
	}

	afterNameParts := strings.Split(afterNamePath, "/")
	category := afterNameParts[1]
	lastPart := afterNameParts[2]

	switch category {
	case "blobs":
		h.blobHandler.GetBlobHandler(c, name, lastPart)
	case "manifests":
		h.manifestHandler.GetManifestHandler(c, name, lastPart)
	case "tags":
		h.manifestHandler.GetTagsHandler(c, name)
	default:
		slog.Error("path is invalid: " + remainPath)
		c.JSON(http.StatusNotFound, "")
	}
}

// func (h FacadeHandler) HandlePOST

// partsNumAfterName とはパスを / で区切った際に name の後ろにあるパートの数
//
// 例）/<name>/blobs/hoge = 2, /<name>/blobs/upload/uuid = 3
func pickUpName(path string, partsNumAfterName int) (name string, afterNamePath string, err error) {
	parts := strings.Split(path, "/")
	partsLen := len(parts)
	if partsLen <= partsNumAfterName {
		err = errors.New("path is invalid")
		return
	}

	for i := 0; i < partsNumAfterName; i++ {
		afterNamePath += ("/" + parts[partsLen-partsNumAfterName+i])
	}

	name = strings.TrimSuffix(path, afterNamePath)
	return
}
