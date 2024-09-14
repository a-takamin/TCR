package persister

import (
	"io"

	"github.com/a-takamin/tcr/internal/dto"
	"github.com/a-takamin/tcr/internal/model"
)

type BlobPersister interface {
	GetBlob(name string, digest string) (model.Blob, error)
	UploadBlob(key string, blob io.ReadCloser) error
	GetChunkedBlobUploadProgress(name string) (dto.BlobUploadProgress, error)
	PutChunkedBlobUpdateProgress(newProgress dto.BlobUploadProgress) error
	PutChunkedBlobConcatenateProgress(concatProgress dto.BlobConcatenateProgress) error
	GetChunkedBlobConcatenateProgress(digest string) (dto.BlobConcatenateProgress, error)
	DeleteBlob(dto.DeleteBlobInput) error
}
