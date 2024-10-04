package apperrors

import (
	"errors"
	"fmt"
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

// TODO: persister をうまくリファクタできたらもうすこし具体的にする
var TCRERR_PERSISTER_ERROR = &TCRError{Message: "永続化層でエラーが発生しました"}

// TODO: これももうちょっと細かくした方が良いかも
var TCRERR_LOGIC_ERROR = &TCRError{Message: "内部でエラーが発生しました"}

var TCRERR_TAG_INVALID = &TCRError{Message: "tag の形式が不正です"}
var TCRERR_NAME_INVALID = &TCRError{Message: "name の形式が不正です"}
var TCRERR_MANIFEST_INVALID = &TCRError{Message: "manifest の形式が不正です"}
var TCRERR_MANIFEST_NOT_FOUND = &TCRError{Message: "対象の manifest がありません"}
var TCRERR_NAME_NOT_FOUND = &TCRError{Message: "対象の name を持つリポジトリがありません"}
var TCRERR_DIGEST_INVALID = &TCRError{Message: "digest の形式が不正です"}
var TCRERR_BLOB_NOT_FOUND = &TCRError{Message: "対象の blob がありません"}
var TCRERR_UNKNOWN = &TCRError{Message: "不明なエラー。このエラーが出た場合は適切な TCRError オブジェクトが利用されるようにエラー処理を修正してください"}

// OCI Error Code はすべてのエラーレスポンスに対して必須というわけではないので、TCR のエラーを作る
type TCRError struct {
	Message string
	Err     error
}

func (e TCRError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Err.Error())
	}
	return e.Message
}

func (e *TCRError) Unwrap() error {
	return e.Err
}

func (e *TCRError) Wrap(err error) error {
	e.Err = err
	return e
}

// TODO: これ作ってるけど、API によってはエラーコードとステータスコードが変わるので、結局使わなくなって API ごとに似たような処理を書くかも
func CreateErrorResponse(err error) (uint, OCIErrorResponse) {
	var tcrErr TCRError
	if !errors.As(err, &tcrErr) {
		err = TCRERR_UNKNOWN.Wrap(err)
	}
	switch err {
	case TCRERR_NAME_INVALID:
		return 400, NAME_INVALID.CreateResponse("")
	default:
		return 500, OCIErrorResponse{Errors: []OCIError{{Detail: "不明なエラー"}}}
	}
}
