package persister

import (
	"github.com/a-takamin/tcr/internal/dto"
)

type BlobPersister interface {
	ExistsBlob(input dto.ExistsBlobInput) (bool, error)
	FindBlob(input dto.FindBlobInput) (dto.FindBlobOutput, error)
	FindChunkedBlob(input dto.FindChunkedBlobInput) (dto.FindBlobOutput, error)
	SaveBlob(input dto.SaveBlobInput) error
	SaveChunkedBlob(input dto.SaveChunkedBlobInput) error
	DeleteBlob(input dto.DeleteBlobInput) error
}
