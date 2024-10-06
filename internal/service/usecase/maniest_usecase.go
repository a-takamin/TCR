package usecase

import (
	"errors"

	"github.com/a-takamin/tcr/internal/apperrors"
	"github.com/a-takamin/tcr/internal/dto"
	"github.com/a-takamin/tcr/internal/interface/persister"
	"github.com/a-takamin/tcr/internal/model"
	"github.com/a-takamin/tcr/internal/service/domain"
)

type ManifestUseCase struct {
	maniRepo persister.ManifestPersister
	repoRepo persister.RepositoryPersister
}

func NewManifestUseCase(maniRepo persister.ManifestPersister, repoRepo persister.RepositoryPersister) *ManifestUseCase {
	return &ManifestUseCase{
		maniRepo: maniRepo,
		repoRepo: repoRepo,
	}
}

func (u ManifestUseCase) ExistsManifest(metadata model.ManifestMetadata) (dto.GetManifestResponse, error) {
	return u.GetManifest(metadata)
}

func (u ManifestUseCase) GetManifest(metadata model.ManifestMetadata) (dto.GetManifestResponse, error) {
	err := domain.ValidateName(metadata.Name)
	if err != nil {
		return dto.GetManifestResponse{}, apperrors.TCRERR_NAME_INVALID
	}

	exists, err := u.maniRepo.ExistsManifest(dto.ExistsManifestInput{
		Name:      metadata.Name,
		Reference: metadata.Reference,
	})
	if err != nil {
		return dto.GetManifestResponse{}, apperrors.TCRERR_PERSISTER_ERROR.Wrap(err)
	}
	if !exists {
		return dto.GetManifestResponse{}, apperrors.TCRERR_NAME_NOT_FOUND
	}

	resp, err := u.maniRepo.FindManifest(dto.FindManifestInput{
		Name:      metadata.Name,
		Reference: metadata.Reference,
	})
	if err != nil {
		return dto.GetManifestResponse{}, apperrors.TCRERR_PERSISTER_ERROR.Wrap(err)
	}

	return dto.GetManifestResponse{
		Manifest: string(resp.Manifest),
		Digest:   resp.Digest,
	}, nil
}

func (u ManifestUseCase) GetTags(name string) (dto.GetTagsResponse, error) {
	err := domain.ValidateName(name)
	if err != nil {
		return dto.GetTagsResponse{}, err
	}

	existsName, err := u.repoRepo.ExistsRepository(dto.ExistsRepositoryInput{
		Name: name,
	})
	if err != nil {
		return dto.GetTagsResponse{}, apperrors.TCRERR_PERSISTER_ERROR.Wrap(err)
	}
	resp, err := u.maniRepo.GetTags(name)
	if err != nil {
		return dto.GetTagsResponse{}, apperrors.TCRERR_PERSISTER_ERROR.Wrap(err)
	}
	if !existsName {
		return dto.GetTagsResponse{}, apperrors.TCRERR_NAME_NOT_FOUND
	}
	return resp, nil
}

func (u ManifestUseCase) PutManifest(metadata model.ManifestMetadata, manifest []byte) error {
	err := domain.ValidateName(metadata.Name)
	if err != nil {
		return apperrors.TCRERR_NAME_INVALID
	}

	err = domain.ValidateManifest(metadata, manifest)
	if err != nil {
		return apperrors.TCRERR_MANIFEST_INVALID.Wrap(err)
	}

	existsName, err := u.repoRepo.ExistsRepository(dto.ExistsRepositoryInput{
		Name: metadata.Name,
	})
	if err != nil {
		return apperrors.TCRERR_PERSISTER_ERROR.Wrap(err)
	}
	if !existsName {
		return apperrors.TCRERR_NAME_NOT_FOUND
	}

	calcdDigest, err := domain.CalcManifestDigestRefactor(manifest)
	if err != nil {
		return err
	}
	isDigest := domain.IsDigest(metadata.Reference)
	var tag string
	if isDigest {
		if calcdDigest != metadata.Reference {
			return errors.New("digest does not match")
		}
		tag = calcdDigest // tag がない場合は digest を tag にする
	} else {
		tag = metadata.Reference
	}

	err = u.maniRepo.SaveManifest(dto.SaveManifestInput{
		Name:     metadata.Name,
		Tag:      tag,
		Digest:   calcdDigest,
		Manifest: manifest,
	})

	if err != nil {
		return errors.New("適切なエラーを設定してください")
	}
	return nil
}

func (u ManifestUseCase) DeleteManifest(metadata model.ManifestMetadata) error {
	err := domain.ValidateName(metadata.Name)
	if err != nil {
		return apperrors.TCRERR_NAME_INVALID
	}

	existsName, err := u.repoRepo.ExistsRepository(dto.ExistsRepositoryInput{
		Name: metadata.Name,
	})
	if err != nil {
		return apperrors.TCRERR_PERSISTER_ERROR.Wrap(err)
	}
	if !existsName {
		return apperrors.TCRERR_NAME_NOT_FOUND
	}

	err = u.maniRepo.DeleteManifest(dto.DeleteManifestInput{
		Name:      metadata.Name,
		Reference: metadata.Reference,
	})

	if err != nil {
		return err
	}
	return nil
}
