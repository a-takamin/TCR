package usecase

import (
	"github.com/a-takamin/tcr/internal/dto"
	"github.com/a-takamin/tcr/internal/interface/persister"
	"github.com/a-takamin/tcr/internal/model"
	"github.com/a-takamin/tcr/internal/service/domain"
)

type ManifestUseCase struct {
	repo persister.ManifestPersister
}

func NewManifestUseCase(repo persister.ManifestPersister) *ManifestUseCase {
	return &ManifestUseCase{
		repo: repo,
	}
}

func (u ManifestUseCase) ExistsManifest(metadata model.ManifestMetadata) (dto.GetManifestResponse, error) {
	return u.GetManifest(metadata)
}

func (u ManifestUseCase) GetManifest(metadata model.ManifestMetadata) (dto.GetManifestResponse, error) {
	err := domain.ValidateNameSpace(metadata.Name)
	if err != nil {
		return dto.GetManifestResponse{}, err
	}

	manifest, err := u.repo.GetManifest(metadata)
	if err != nil {
		return dto.GetManifestResponse{}, err
	}

	digest, err := domain.CalcManifestDigest(manifest)
	if err != nil {
		return dto.GetManifestResponse{}, err
	}

	return dto.GetManifestResponse{
		Manifest: manifest,
		Digest:   digest,
	}, nil
}

func (u ManifestUseCase) PutManifest(metadata model.ManifestMetadata, manifest model.Manifest) error {
	err := domain.ValidateNameSpace(metadata.Name)
	if err != nil {
		return err
	}

	err = u.repo.PutManifest(metadata, manifest)
	if err != nil {
		return err
	}
	return nil
}

func (u ManifestUseCase) DeleteManifest(metadata model.ManifestMetadata) error {
	err := domain.ValidateNameSpace(metadata.Name)
	if err != nil {
		return err
	}

	return u.repo.DeleteManifest(metadata)
}

func (u ManifestUseCase) GetTags(name string) (dto.GetTagsResponse, error) {
	err := domain.ValidateNameSpace(name)
	if err != nil {
		return dto.GetTagsResponse{}, err
	}
	return u.repo.GetTags(name)
}
