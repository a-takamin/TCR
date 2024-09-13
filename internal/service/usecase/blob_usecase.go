package usecase

import (
	"errors"
	"fmt"

	"github.com/a-takamin/tcr/internal/apperrors"
	"github.com/a-takamin/tcr/internal/dto"
	"github.com/a-takamin/tcr/internal/interface/persister"
	"github.com/a-takamin/tcr/internal/model"
	"github.com/a-takamin/tcr/internal/service/domain"
	"github.com/google/uuid"
)

type BlobUseCase struct {
	blob *domain.BlobDomain
	repo persister.BlobPersister
}

func NewBlobUseCase(s *domain.BlobDomain, r persister.BlobPersister) *BlobUseCase {
	return &BlobUseCase{
		blob: s,
		repo: r,
	}
}

func (u BlobUseCase) GetBlob(m dto.GetBlobInput) (model.Blob, error) {
	// Get はほぼやることない
	err := u.blob.ValidateNameSpace(m.Name)
	if err != nil {
		return model.Blob{}, err
	}
	blob, err := u.repo.GetBlob(m.Name, m.Digest)
	if err != nil {
		return model.Blob{}, err
	}
	return blob, err
}

func (u BlobUseCase) StartBlobUpload(name string) (string, error) {
	uid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("/v2/%s/blobs/uploads/%s", name, uid), nil
}

func (u BlobUseCase) UploadMonolithicBlob(input dto.UploadMonolithicBlobInput) error {
	err := u.blob.ValidateNameSpace(input.Name)
	if err != nil {
		return err
	}
	err = u.blob.ValidateDigest(input.Digest)
	if err != nil {
		return err
	}
	err = u.repo.UploadBlob(input.Digest, input.Blob)
	if err != nil {
		return err
	}
	return nil
}

// int64: アップロードに成功したバイト数
//
// error: エラー
func (u BlobUseCase) UploadChunkedBlob(input dto.UploadChunkedBlobInput) (int64, error) {
	err := u.blob.ValidateNameSpace(input.Name)
	if err != nil {
		return 0, err
	}
	err = u.blob.ValidateContentRange(input.ContentRange)
	if err != nil {
		return 0, err
	}
	startByte, err := u.blob.GetContentRangeStart(input.ContentRange)
	if err != nil {
		return 0, err
	}
	info, err := u.repo.GetChunkedBlobUploadProgress(input.Uuid)
	if err != nil {
		return 0, err
	}
	if info.Done {
		return info.ByteUploaded, apperrors.ErrAllChunksAreAlreadyUploaded
	}
	if startByte != info.ByteUploaded {
		return info.ByteUploaded, apperrors.ErrChunkIsNotInSequence
	}

	input.Key = fmt.Sprintf("/chunk/%s/%d", input.Uuid, info.NextChunkNo)
	err = u.repo.UploadBlob(input.Key, input.Blob)
	if err != nil {
		return info.ByteUploaded, err
	}

	info.Uuid = input.Uuid
	info.NextChunkNo += 1
	info.ByteUploaded += input.ContentLength
	info.Done = input.IsLast

	err = u.repo.PutChunkedBlobUpdateProgress(info)
	if err != nil {
		return (info.ByteUploaded - input.ContentLength), err
	}
	return info.ByteUploaded, nil
}

func (u BlobUseCase) UploadLastChunkedBlob(input dto.UploadChunkedBlobInput) (int64, error) {
	offset, err := u.UploadChunkedBlob(input)
	if err != nil && !errors.Is(err, apperrors.ErrAllChunksAreAlreadyUploaded) {
		return offset, err
	}

	// temp にアップロードされた過去のレイヤーを結合する処理を非同期で実施。それをトリガー。今回は SQS
	// 非同期の処理のステータスを確認できるテーブルに、 uuid を持つ temp たちが in progress であることを挿入。digest も。
	err = u.StartBlobConcat(input.Digest)
	if err != nil {
		return offset, err
	}

	concatInput := dto.BlobConcatenateProgress{
		Digest: input.Digest,
		Status: "doing",
	}

	err = u.repo.PutChunkedBlobConcatenateProgress(concatInput)

	if err != nil {
		return offset, err
	}

	return offset, nil
}
func (u BlobUseCase) StartBlobConcat(digest string) error {
	// TODO: SQS を呼ぶ
	return nil
}
