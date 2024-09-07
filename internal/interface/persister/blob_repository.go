package persister

import (
	"io"

	"github.com/a-takamin/tcr/internal/model"
)

type BlobPersister interface {
	GetBlob(metadata model.BlobMetadata) (model.Blob, error)
	UploadBlob(metadata model.BlobUploadMetadata, blob io.Reader) error
	GetChunkedBlobUploadProgress(name string) (model.BlobUploadProgress, error)
	PutChunkedBlobUpdateProgress(newProgress model.BlobUploadProgress) error
	// PutBlob(metadata model.BlobMetadata, manifest model.Blob) error
	// DeleteBlob(metadata model.BlobMetadata) error
}
