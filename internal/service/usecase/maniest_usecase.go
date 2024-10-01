package usecase

import (
	"bytes"
	"encoding/base64"
	"encoding/json"

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

	var out bytes.Buffer
	json.Indent(&out, []byte(manifest), "", "\t")
	b := out.Bytes()
	if err != nil {
		return dto.GetManifestResponse{}, err
	}
	digest, err := domain.CalcManifestDigest(b)
	if err != nil {
		return dto.GetManifestResponse{}, err
	}

	return dto.GetManifestResponse{
		Manifest: manifest,
		Digest:   digest,
	}, nil
}

func (u ManifestUseCase) PutManifest(metadata model.ManifestMetadata, manifest []byte) error {
	err := domain.ValidateNameSpace(metadata.Name)
	if err != nil {
		return err
	}

	err = domain.ValidateManifest(metadata, manifest)
	if err != nil {
		return err
	}

	encodedManifest := base64.StdEncoding.EncodeToString(manifest)

	// string がいくべきか？
	err = u.repo.PutManifest(metadata, encodedManifest)
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
