package service

import (
	"errors"
	"fmt"
	"io"

	"github.com/a-takamin/tcr/internal/interface/persister"
	"github.com/a-takamin/tcr/internal/model"
	"github.com/a-takamin/tcr/internal/service/utils"
	"github.com/google/uuid"
)

type BlobService struct {
	repo persister.BlobPersister
}

func NewBlobService(repository persister.BlobPersister) *BlobService {
	return &BlobService{repo: repository}
}

func (s BlobService) GetBlob(metadata model.BlobMetadata) (model.Blob, error) {
	err := utils.ValidateName(metadata.Name)
	if err != nil {
		return model.Blob{}, err
	}

	blob, err := s.repo.GetBlob(metadata)
	if err != nil {
		return model.Blob{}, err
	}

	// if c.Request.Header["Content-Type"][0] != blob.MediaType {
	// 	return model.Blob{}, errors.New("Content-Type does not match")
	// }

	return blob, nil
}

func (s BlobService) StartBlobUpload(name string) (string, error) {
	uid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("/v2/%s/blobs/uploads/%s", name, uid), nil
}

func (s BlobService) UploadBlob() error {
	return nil
}

func (s BlobService) UploadBlobMonolithically(metadata model.BlobUploadMetadata, blob io.ReadCloser) error {
	err := utils.ValidateName(metadata.Name)
	if err != nil {
		return err
	}
	if !utils.IsDigest(metadata.Digest) {
		return errors.New("digest is invalid")
	}

	metadata.Key = s.CreateKey(metadata)

	err = s.repo.UploadBlob(metadata, blob)
	if err != nil {
		return err
	}
	return nil
}

func (s BlobService) CreateKey(metadata model.BlobUploadMetadata) string {
	if metadata.IsChunkUpload {
		// 次のチャンクの順番IDを取得してセット
		var chunkId string
		return fmt.Sprintf("/chunk/%s/%s", metadata.Uuid, chunkId)

	}
	return metadata.Uuid
}

// func (s BlobService) PutBlob(metadata model.BlobMetadata, blob model.Blob) error {
// 	err := utils.ValidateName(metadata.Name)
// 	if err != nil {
// 		return err
// 	}

// 	err = s.repo.PutBlob(metadata, blob)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (s BlobService) DeleteBlob(metadata model.BlobMetadata) error {
// 	err := utils.ValidateName(metadata.Name)
// 	if err != nil {
// 		return err
// 	}

// 	err = s.repo.DeleteBlob(metadata)
// 	if err != nil {
// 		return err
// 	}
// 	return nil

// }
