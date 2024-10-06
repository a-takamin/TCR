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
	blobRepo     persister.BlobPersister
	progressRepo persister.BlobUploadProgressPersister
	repoRepo     persister.RepositoryPersister
}

func NewBlobUseCase(blobRepo persister.BlobPersister, progressRepo persister.BlobUploadProgressPersister, repoRepo persister.RepositoryPersister) *BlobUseCase {
	return &BlobUseCase{
		blobRepo:     blobRepo,
		progressRepo: progressRepo,
		repoRepo:     repoRepo,
	}
}

func (u BlobUseCase) ExistsBlob(input dto.FindBlobInput) (model.Blob, error) {
	return u.GetBlob(input)
}

func (u BlobUseCase) GetBlob(input dto.FindBlobInput) (model.Blob, error) {
	err := domain.ValidateName(input.Name)
	if err != nil {
		return model.Blob{}, apperrors.TCRERR_NAME_INVALID
	}
	err = domain.ValidateDigest(input.Digest)
	if err != nil {
		return model.Blob{}, apperrors.TCRERR_DIGEST_INVALID
	}
	existsName, err := u.repoRepo.ExistsRepository(dto.ExistsRepositoryInput{
		Name: input.Name,
	})
	if err != nil {
		return model.Blob{}, apperrors.TCRERR_PERSISTER_ERROR.Wrap(err)
	}
	if !existsName {
		return model.Blob{}, apperrors.TCRERR_NAME_NOT_FOUND
	}
	existsBlob, err := u.blobRepo.ExistsBlob(dto.ExistsBlobInput{
		Name:   input.Name,
		Digest: input.Digest,
	})
	if err != nil {
		return model.Blob{}, apperrors.TCRERR_PERSISTER_ERROR.Wrap(err)
	}
	if !existsBlob {
		return model.Blob{}, apperrors.TCRERR_BLOB_NOT_FOUND
	}

	resp, err := u.blobRepo.FindBlob(dto.FindBlobInput{
		Name:   input.Name,
		Digest: input.Digest,
	})
	if err != nil {
		return model.Blob{}, apperrors.TCRERR_PERSISTER_ERROR.Wrap(err)
	}
	return model.Blob{
		Name:   input.Name,
		Digest: input.Digest,
		Blob:   resp.Blob,
	}, nil
}

func (u BlobUseCase) StartBlobUpload(name string) (string, error) {
	uid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	err = u.progressRepo.SaveBlobUploadProgress(dto.SaveBlobUploadProgressInput{
		Uuid:         uid.String(),
		NextChunkNo:  0,
		ByteUploaded: 0,
		Digest:       "",
	})
	if err != nil {
		return "", err
	}
	err = u.repoRepo.SaveRepository(dto.SaveRepositoryInput{
		Name: name,
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("/v2/%s/blobs/uploads/%s", name, uid), nil
}

func (u BlobUseCase) UploadMonolithicBlob(input dto.UploadMonolithicBlobInput) error {
	err := domain.ValidateName(input.Name)
	if err != nil {
		return err
	}
	err = domain.ValidateDigest(input.Digest)
	if err != nil {
		return err
	}

	err = u.blobRepo.SaveBlob(dto.SaveBlobInput{
		Name:   input.Name,
		Digest: input.Digest,
		Blob:   input.Blob,
	})
	if err != nil {
		return err
	}
	return nil
}

// int64: アップロードに成功したバイト数
//
// error: エラー
func (u BlobUseCase) UploadChunkedBlob(input dto.UploadChunkedBlobInput) (int64, error) {
	err := domain.ValidateName(input.Name)
	if err != nil {
		return 0, err
	}
	err = domain.ValidateContentRange(input.ContentRange)
	if err != nil {
		return 0, err
	}
	startByte, err := domain.GetContentRangeStart(input.ContentRange)
	if err != nil {
		return 0, err
	}
	endByte, err := domain.GetContentRangeEnd(input.ContentRange)
	if err != nil {
		return 0, err
	}

	info, err := u.progressRepo.FindBlobUploadProgress(dto.FindBlobUploadProgressInput{
		Uuid: input.Uuid,
	})
	if err != nil {
		return 0, err
	}

	// TODO: 綺麗にする
	// if info.ByteUploaded == 0 {
	if startByte != info.ByteUploaded {
		return info.ByteUploaded, apperrors.ErrChunkIsNotInSequence
	}
	// } else {
	// 	if startByte != info.ByteUploaded+1 {
	// 		return info.ByteUploaded, apperrors.ErrChunkIsNotInSequence
	// 	}
	// }

	err = u.blobRepo.SaveChunkedBlob(dto.SaveChunkedBlobInput{
		Name:       input.Name,
		Uuid:       input.Uuid,
		ChunkSeqNo: info.NextChunkNo,
		Blob:       input.Blob,
	})
	if err != nil {
		return info.ByteUploaded, err
	}

	err = u.progressRepo.SaveBlobUploadProgress(dto.SaveBlobUploadProgressInput{
		Uuid:         input.Uuid,
		ByteUploaded: info.ByteUploaded + input.ContentLength,
		NextChunkNo:  info.NextChunkNo + 1,
		Digest:       input.Digest,
	})
	if err != nil {
		return info.ByteUploaded, err
	}
	return endByte, nil
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

	info, err := u.progressRepo.FindBlobUploadProgress(dto.FindBlobUploadProgressInput{
		Uuid: input.Uuid,
	})
	if err != nil {
		return offset, err
	}

	err = u.progressRepo.SaveBlobUploadProgress(dto.SaveBlobUploadProgressInput{
		Uuid:         info.Uuid,
		ByteUploaded: info.ByteUploaded,
		NextChunkNo:  info.NextChunkNo,
		Digest:       input.Digest, // Digest を登録
	})
	if err != nil {
		return offset, err
	}

	err = u.StartBlobConcat(input.Name, input.Uuid, input.Digest)
	if err != nil {
		return offset, err
	}

	return offset, nil
}
func (u BlobUseCase) StartBlobConcat(name string, uuid string, digest string) error {
	// TODO: 非同期でやりたい
	// TODO: ストリームでやりたい。今のままでは巨大なイメージに押しつぶされる

	info, err := u.progressRepo.FindBlobUploadProgress(dto.FindBlobUploadProgressInput{
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	if info.Digest == "" {
		// TODO: エラー処理。今は続ける。
		slog.Warn("Digest does not exist")
	}

	chunkNums := info.NextChunkNo
	if chunkNums <= 0 {
		return errors.New("NextChunkNo should be greater than or equal to 0")
	}
	var concatBlob []byte
	for i := 0; i != chunkNums; i++ {
		resp, err := u.blobRepo.FindChunkedBlob(dto.FindChunkedBlobInput{
			Name:       name,
			Uuid:       uuid,
			ChunkSeqNo: i,
		})
		if err != nil {
			// TODO: ちゃんとエラーハンドリング
			slog.Warn("chunk not found")
			return apperrors.ErrBlobNotFound
		}
		concatBlob = append(concatBlob, resp.Blob...)
	}

	u.blobRepo.SaveBlob(dto.SaveBlobInput{
		Name:   name,
		Digest: digest,
		Blob:   bytes.NewReader(concatBlob),
	})
	if err != nil {
		return err
	}
	return nil
}

func (u BlobUseCase) DeleteBlob(input dto.DeleteBlobInput) error {
	_, err := u.ExistsBlob(dto.FindBlobInput{
		Name:   input.Name,
		Digest: input.Digest,
	})
	if err != nil {
		return err
	}
	err = domain.ValidateName(input.Name)
	if err != nil {
		return err
	}
	return u.blobRepo.DeleteBlob(input)
}

// TODO: モノリスかラストチャンクかの見分けをもう少しちゃんと考える
func (u BlobUseCase) IsChunkedUpload(name string, uuid string) (bool, error) {
	info, err := u.progressRepo.FindBlobUploadProgress(dto.FindBlobUploadProgressInput{
		Uuid: uuid,
	})

	if err != nil {
		return false, err
	}
	if info.ByteUploaded > 0 {
		return true, nil
	}
	if info.Digest != "" {
		return true, nil
	}
	if info.NextChunkNo > 0 {
		return true, nil
	}
	return false, nil
}

func (u BlobUseCase) GetBlobUploadOffset(name string, uuid string) (int64, error) {
	info, err := u.progressRepo.FindBlobUploadProgress(dto.FindBlobUploadProgressInput{
		Uuid: uuid,
	})
	if err != nil {
		return 0, err
	}
	return info.ByteUploaded - 1, nil
}
