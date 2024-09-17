package apperrors

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO: 実はOCIで定義されている。作り直し。
var ErrInvalidName = errors.New("name is invalid")
var ErrInvalidManifest = errors.New("manifest is invalid")
var ErrInvalidReference = errors.New("reference is invalid")
var ErrManifestNotFound = errors.New("manifest not found")
var ErrBlobNotFound = errors.New("blob not found")
var ErrInvalidContentRange = errors.New("Content-Range format is invalid")
var ErrChunkIsNotInSequence = errors.New("chunk is not in sequence")
var ErrAllChunksAreAlreadyUploaded = errors.New("all chunks are already uploaded")

func ErrorHanlder(c *gin.Context, err error) {
	switch err {
	case ErrInvalidName:
		c.JSON(http.StatusBadRequest, err)
	case ErrInvalidReference:
		c.JSON(http.StatusBadRequest, err)
	case ErrManifestNotFound:
		c.JSON(http.StatusNotFound, err)
	case ErrBlobNotFound:
		c.JSON(http.StatusNotFound, err)
	default:
		c.JSON(http.StatusInternalServerError, err)
	}
}
