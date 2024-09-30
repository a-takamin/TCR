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

var BLOB_UNKNOWN = &TCRError{ErrorCode: "BLOB_UNKNOWN", ErrorMessage: "blob unknown to registry"}
var BLOB_UPLOAD_INVALID = &TCRError{ErrorCode: "BLOB_UPLOAD_INVALID", ErrorMessage: "blob upload invalid"}
var BLOB_UPLOAD_UNKNOWN = &TCRError{ErrorCode: "BLOB_UPLOAD_UNKNOWN", ErrorMessage: "blob upload unknown to registry"}
var DIGEST_INVALID = &TCRError{ErrorCode: "DIGEST_INVALID", ErrorMessage: "provided digest did not match uploaded content"}
var MANIFEST_BLOB_UNKNOWN = &TCRError{ErrorCode: "MANIFEST_BLOB_UNKNOWN", ErrorMessage: "blob unknown to registry"}
var MANIFEST_INVALID = &TCRError{ErrorCode: "MANIFEST_INVALID", ErrorMessage: "manifest invalid"}
var MANIFEST_UNKNOWN = &TCRError{ErrorCode: "MANIFEST_UNKNOWN", ErrorMessage: "manifest unknown"}
var MANIFEST_UNVERIFIED = &TCRError{ErrorCode: "MANIFEST_UNVERIFIED", ErrorMessage: "manifest failed signature verification"}
var NAME_INVALID = &TCRError{ErrorCode: "NAME_INVALID", ErrorMessage: "invalid repository name"}
var NAME_UNKNOWN = &TCRError{ErrorCode: "NAME_UNKNOWN", ErrorMessage: "repository name not known to registry"}
var SIZE_INVALID = &TCRError{ErrorCode: "SIZE_INVALID", ErrorMessage: "provided length did not match content length"}
var TAG_INVALID = &TCRError{ErrorCode: "TAG_INVALID", ErrorMessage: "manifest tag did not match URI"}
var UNAUTHORIZED = &TCRError{ErrorCode: "UNAUTHORIZED", ErrorMessage: "authentication required"}
var DENIED = &TCRError{ErrorCode: "DENIED", ErrorMessage: "requested access to the resource is denied"}
var UNSUPPORTED = &TCRError{ErrorCode: "UNSUPPORTED", ErrorMessage: "The operation is unsupported"}
var INTERNAL_SERVER_ERROR = &TCRError{ErrorCode: "INTERNAL_SERVER_ERROR", ErrorMessage: "internal server error"}

type TCRError struct {
	ErrorCode    string
	ErrorMessage string
	Err          error `json:"-"` // 内部的なエラーはユーザーに見せない
}

func (e TCRError) Error() string {
	return e.Err.Error()
}

func (e *TCRError) Unwrap() error {
	return e.Err
}

func (e *TCRError) Wrap(err error) error {
	e.Err = err
	return e
}

// TODO: 直す
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
