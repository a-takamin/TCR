package persister

import "github.com/a-takamin/tcr/internal/dto"

type BlobUploadProgressPersister interface {
	FindBlobUploadProgress(input dto.FindBlobUploadProgressInput) (dto.FindBlobUploadProgressOutput, error)
	SaveBlobUploadProgress(input dto.SaveBlobUploadProgressInput) error
}
