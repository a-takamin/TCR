package usecase

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"

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

func (u BlobUseCase) ExistsBlob(input dto.GetBlobInput) (model.Blob, error) {
	return u.GetBlob(input)
}

func (u BlobUseCase) GetBlob(input dto.GetBlobInput) (model.Blob, error) {
	err := u.blob.ValidateNameSpace(input.Name)
	if err != nil {
		return model.Blob{}, err
	}
	blob, err := u.repo.GetBlob(input.Name, input.Digest)
	if err != nil {
		// TODO: ちゃんとエラーハンドリング
		return model.Blob{}, apperrors.ErrBlobNotFound
	}
	return blob, err
}

func (u BlobUseCase) StartBlobUpload(name string) (string, error) {
	uid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	info := dto.BlobUploadProgress{
		Uuid:         uid.String(),
		NextChunkNo:  0,
		ByteUploaded: 0,
		Done:         false,
	}
	err = u.repo.PutChunkedBlobUpdateProgress(info)
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
	key := fmt.Sprintf("%s/%s", input.Name, input.Digest)
	err = u.repo.UploadBlob(key, input.Blob)
	if err != nil {
		return err
	}
	return nil
}

// int64: アップロードに成功したバイト数
//
// error: エラー
func (u BlobUseCase) UploadChunkedBlob(input dto.UploadChunkedBlobInput) (int64, error) {
	// err := u.blob.ValidateNameSpace(input.Name)
	// if err != nil {
	// 	return 0, err
	// }
	// err = u.blob.ValidateContentRange(input.ContentRange)
	// if err != nil {
	// 	return 0, err
	// }
	// startByte, err := u.blob.GetContentRangeStart(input.ContentRange)
	// if err != nil {
	// 	return 0, err
	// }
	info, err := u.repo.GetChunkedBlobUploadProgress(input.Uuid)
	if err != nil {
		return 0, err
	}
	// if info.Done {
	// 	return info.ByteUploaded, apperrors.ErrAllChunksAreAlreadyUploaded
	// }
	// // TODO: 綺麗にする
	// if info.ByteUploaded == 0 {
	// 	if startByte != info.ByteUploaded {
	// 		return info.ByteUploaded, apperrors.ErrChunkIsNotInSequence
	// 	}
	// } else {
	// 	if startByte != info.ByteUploaded+1 {
	// 		return info.ByteUploaded, apperrors.ErrChunkIsNotInSequence
	// 	}
	// }

	input.Key = fmt.Sprintf("/%s/chunk/%s/%d", input.Name, input.Uuid, info.NextChunkNo)
	err = u.repo.UploadBlob(input.Key, input.Blob)
	if err != nil {
		return info.ByteUploaded, err
	}

	info.Uuid = input.Uuid
	info.NextChunkNo += 1
	info.ByteUploaded += input.ContentLength - 1
	info.Done = input.IsLast

	err = u.repo.PutChunkedBlobUpdateProgress(info)
	if err != nil {
		return (info.ByteUploaded - input.ContentLength), err
	}
	return info.ByteUploaded, nil
}

func (u BlobUseCase) UploadLastChunkedBlob(input dto.UploadChunkedBlobInput) (int64, error) {
	var offset int64
	var err error

	if input.ContentLength != 0 {
		// Last Upload with Blob
		offset, err = u.UploadChunkedBlob(input)
		if err != nil && !errors.Is(err, apperrors.ErrAllChunksAreAlreadyUploaded) {
			return offset, err
		}
	}

	// temp にアップロードされた過去のレイヤーを結合する処理を非同期で実施。それをトリガー。今回は SQS
	// 非同期の処理のステータスを確認できるテーブルに、 uuid を持つ temp たちが in progress であることを挿入。digest も。
	err = u.StartBlobConcat(input.Name, input.Uuid, input.Digest)
	if err != nil {
		return offset, err
	}

	info, err := u.repo.GetChunkedBlobUploadProgress(input.Uuid)
	if err != nil {
		return 0, err
	}

	info.Done = true
	err = u.repo.PutChunkedBlobUpdateProgress(info)

	if err != nil {
		return offset, err
	}

	return offset, nil
}
func (u BlobUseCase) StartBlobConcat(name string, uuid string, digest string) error {
	// TODO: 非同期でやりたい
	// TODO: ストリームでやりたい。今のままでは巨大なイメージに押しつぶされる

	info, err := u.repo.GetChunkedBlobUploadProgress(uuid)
	if err != nil {
		return err
	}
	if !info.Done {
		// TODO: エラー処理。今は続ける。
		slog.Warn("Done is false")
	}

	chunkNums := info.NextChunkNo
	if chunkNums < 0 {
		return errors.New("NextChunkNo should be greater than or equal to 0")
	}
	var concatBlob []byte
	for i := 0; i != chunkNums; i++ {
		key := fmt.Sprintf("/%s/chunk/%s/%d", name, uuid, i)
		// TODO: 永続化層にロジックを持たせているせいで苦労している。直す。
		blobModel, err := u.repo.GetBlob("", key)
		if err != nil {
			// TODO: ちゃんとエラーハンドリング
			slog.Warn("chunk not found")
			return apperrors.ErrBlobNotFound
		}
		concatBlob = append(concatBlob, blobModel.Blob...)
	}
	key := fmt.Sprintf("%s/%s", name, digest)
	err = u.repo.UploadBlob(key, bytes.NewReader(concatBlob))
	if err != nil {
		return err
	}
	return nil
}

func (u BlobUseCase) DeleteBlob(input dto.DeleteBlobInput) error {
	_, err := u.ExistsBlob(dto.GetBlobInput{
		Name:   input.Name,
		Digest: input.Digest,
	})
	if err != nil {
		return err
	}
	err = u.blob.ValidateNameSpace(input.Name)
	if err != nil {
		return err
	}
	return u.repo.DeleteBlob(input)
}

// TODO: モノリスかチャンクかの見分けをもう少しちゃんと考える
func (u BlobUseCase) IsChunkedUpload(name string, uuid string) (bool, error) {
	info, err := u.repo.GetChunkedBlobUploadProgress(uuid)
	if err != nil {
		return false, err
	}
	if info.ByteUploaded > 0 {
		return true, nil
	}
	if info.Done {
		return true, nil
	}
	if info.NextChunkNo > 0 {
		return true, nil
	}
	return false, nil
}

func (u BlobUseCase) GetBlobUploadStatus(name string, uuid string) (int64, error) {
	info, err := u.repo.GetChunkedBlobUploadProgress(uuid)
	if err != nil {
		return 0, err
	}
	return info.ByteUploaded, nil
}
