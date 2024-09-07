package service

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/a-takamin/tcr/apperrors"
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

func (s BlobService) UploadBlob(metadata model.BlobUploadMetadata, blob io.ReadCloser) error {
	err := utils.ValidateName(metadata.Name)
	if err != nil {
		return err
	}

	if metadata.IsChunkUpload {
		info, err := s.repo.GetChunkedBlobUploadProgress(metadata.Name)
		if err != nil {
			return err
		}
		ranges := strings.Split(metadata.ContentRange, "-")
		if len(ranges) != 2 {
			return errors.New("Content-Range format is invalid")
		}
		if ranges[0] != string(info.ByteUploaded) {
			return apperrors.ErrChunkIsNotInSequence
		}
		metadata.Key = fmt.Sprintf("/chunk/%s/%s", metadata.Uuid, string(info.NextChunkNo))
		err = s.repo.UploadBlob(metadata, blob)
		if err != nil {
			return err
		}
		// update
		info.NextChunkNo++
		info.ByteUploaded += metadata.ContentLength
		info.Done = false
		err = s.repo.PutChunkedBlobUpdateProgress(info)
		if err != nil {
			return err
		}
	}

	// Monolithic upload must contain digest
	if !utils.IsDigest(metadata.Digest) {
		return errors.New("digest is invalid")
	}

	metadata.Key = metadata.Uuid

	err = s.repo.UploadBlob(metadata, blob)
	if err != nil {
		return err
	}

	return nil
}

func (s BlobService) UploadChunkedBlob(metadata model.BlobUploadMetadata, blob io.ReadCloser) (int64, error) {
	err := utils.ValidateName(metadata.Name)
	if err != nil {
		return 0, err
	}

	info, err := s.repo.GetChunkedBlobUploadProgress(metadata.Name)
	if err != nil {
		return 0, err
	}
	ranges := strings.Split(metadata.ContentRange, "-")
	if len(ranges) != 2 {
		return info.ByteUploaded, errors.New("Content-Range format is invalid")
	}

	if ranges[0] != strconv.FormatInt(info.ByteUploaded, 10) {
		return info.ByteUploaded, apperrors.ErrChunkIsNotInSequence
	}
	metadata.Key = fmt.Sprintf("/chunk/%s/%d", metadata.Uuid, info.NextChunkNo)
	err = s.repo.UploadBlob(metadata, blob)
	if err != nil {
		return info.ByteUploaded, err
	}
	// update
	info.Name = metadata.Name
	info.NextChunkNo++
	info.ByteUploaded += metadata.ContentLength
	info.Done = false
	err = s.repo.PutChunkedBlobUpdateProgress(info)
	if err != nil {
		// 今回の分を戻す
		return info.ByteUploaded - metadata.ContentLength, err
	}
	return info.ByteUploaded, err
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
