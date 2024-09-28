package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// Gin ではパスの変数にスラッシュを使えないために設けられたハンドラー
type FacadeHandler struct {
	blobHandler     *BlobHandler
	manifestHandler *ManifestHandler
}

func NewFacadeHandler(mh *ManifestHandler, bh *BlobHandler) *FacadeHandler {
	return &FacadeHandler{
		blobHandler:     bh,
		manifestHandler: mh,
	}
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
// "/v2/"
//
// "/v2/:name/blobs/:digest"
//
// "/v2/:name/manifests/:reference"
//
// "/v2/:name/tags/list"
//
// "/v2/:name/blobs/uploads/:uuid"
func (h FacadeHandler) HandleGET(c *gin.Context) {
	remainPath := c.Param("remain")

	// GET /v2/
	if remainPath == "/" {
		c.JSON(http.StatusOK, "")
		return
	}

	// TODO: パスを判断する関数を作る
	// 仕様に載っていない /v2/:name/blobs/uploads/:uuid のおかげで if が生えたため。これを機に綺麗にする
	matched, _ := regexp.MatchString(`/blobs/uploads/`, remainPath)
	if matched {
		name, afterNamePath, err := pickUpName(remainPath, 3)
		if err != nil {
			slog.Error(err.Error())
			c.JSON(http.StatusNotFound, "")
			return
		}
		afterNameParts := strings.Split(afterNamePath, "/")
		uuid := afterNameParts[3]
		h.blobHandler.GetUploadStatusHandler(c, name, uuid)
		return
	}

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

// 以下の API
//
// "/v2/:name/blobs/uploads"
func (h FacadeHandler) HandlePOST(c *gin.Context) {
	remainPath := c.Param("remain")
	name, afterNamePath, err := pickUpName(remainPath, 2)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusNotFound, "")
		return
	}

	afterNameParts := strings.Split(afterNamePath, "/")
	category := afterNameParts[1]

	if category != "blobs" {
		c.JSON(http.StatusNotFound, "")
		return
	}
	h.blobHandler.StartUploadBlobHandler(c, name)
}

// 以下の API
//
// "/v2/:name/manifests/:reference"
//
// "/v2/:name/blobs/uploads/:uuid"
func (h FacadeHandler) HandlePUT(c *gin.Context) {
	remainPath := c.Param("remain")
	matched, _ := regexp.MatchString(`/blobs/uploads/`, remainPath)
	if matched {
		name, afterNamePath, err := pickUpName(remainPath, 3)
		if err != nil {
			slog.Error(err.Error())
			c.JSON(http.StatusNotFound, "")
			return
		}
		afterNameParts := strings.Split(afterNamePath, "/")
		uuid := afterNameParts[3]
		h.blobHandler.UploadBlobHandler(c, name, uuid)
		return
	}

	name, afterNamePath, err := pickUpName(remainPath, 2)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusNotFound, "")
		return
	}
	afterNameParts := strings.Split(afterNamePath, "/")
	category := afterNameParts[1]
	reference := afterNameParts[2]
	if category != "manifests" {
		c.JSON(http.StatusNotFound, "")
	}
	h.manifestHandler.PutManifestHandler(c, name, reference)
}

// 以下の API
//
// "/v2/:name/blobs/uploads/:uuid"
func (h FacadeHandler) HandlePATCH(c *gin.Context) {
	remainPath := c.Param("remain")
	name, afterNamePath, err := pickUpName(remainPath, 3)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusNotFound, "")
		return
	}

	afterNameParts := strings.Split(afterNamePath, "/")
	firstPath := afterNameParts[1]
	secondPath := afterNameParts[2]
	uuid := afterNameParts[3]

	if firstPath != "blobs" || secondPath != "uploads" {
		c.JSON(http.StatusNotFound, "")
		return
	}
	h.blobHandler.UploadChunkedBlobHandler(c, name, uuid)
}

// 以下の API
//
// "/v2/:name/manifests/:reference"
//
// /v2/:name/blobs/:reference
func (h FacadeHandler) HandleDELETE(c *gin.Context) {
	remainPath := c.Param("remain")
	name, afterNamePath, err := pickUpName(remainPath, 2)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusNotFound, "")
		return
	}

	afterNameParts := strings.Split(afterNamePath, "/")
	category := afterNameParts[1]
	reference := afterNameParts[2]

	switch category {
	case "blobs":
		h.blobHandler.DeleteBlobHandler(c, name, reference)
	case "manifests":
		h.manifestHandler.DeleteManifestHandler(c, name, reference)
	default:
		slog.Error("path is invalid: " + remainPath)
		c.JSON(http.StatusNotFound, "")
	}
}

// path から name 部分とその後ろ部分を抜き出す
//
// partsNumAfterName とはパスを / で区切った際に name の後ろにあるパートの数
//
// 例）/<name>/blobs/hoge = 2, /<name>/blobs/upload/uuid = 3
func pickUpName(path string, partsNumAfterName int) (name string, afterNamePath string, err error) {
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

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
