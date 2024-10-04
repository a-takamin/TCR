package usecase

import (
	"bytes"
	"encoding/base64"
	"encoding/json"

	"github.com/a-takamin/tcr/internal/apperrors"
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
		return dto.GetManifestResponse{}, apperrors.TCRERR_NAME_INVALID
	}

	existsName, err := u.repo.ExistsName(metadata.Name)
	if err != nil {
		return dto.GetManifestResponse{}, apperrors.TCRERR_PERSISTER_ERROR.Wrap(err)
	}
	if !existsName {
		return dto.GetManifestResponse{}, apperrors.TCRERR_NAME_NOT_FOUND
	}

	var manifest string
	if domain.IsDigest(metadata.Reference) {
		manifest, err = u.repo.GetManifestByDigest(metadata)
	} else {
		// = tag
		manifest, err = u.repo.GetManifestByTag(metadata)
	}
	if err != nil {
		return dto.GetManifestResponse{}, apperrors.TCRERR_PERSISTER_ERROR.Wrap(err)
	}

	var out bytes.Buffer
	json.Indent(&out, []byte(manifest), "", "\t")
	b := out.Bytes()
	if err != nil {
		return dto.GetManifestResponse{}, apperrors.TCRERR_LOGIC_ERROR.Wrap(err)
	}
	digest, err := domain.CalcManifestDigest(b)
	if err != nil {
		return dto.GetManifestResponse{}, apperrors.TCRERR_LOGIC_ERROR.Wrap(err)
	}

	return dto.GetManifestResponse{
		Manifest: manifest,
		Digest:   digest,
	}, nil
}

func (u ManifestUseCase) PutManifest(metadata model.ManifestMetadata, manifest []byte) error {
	err := domain.ValidateNameSpace(metadata.Name)
	if err != nil {
		return apperrors.TCRERR_NAME_INVALID
	}

	err = domain.ValidateManifest(metadata, manifest)
	if err != nil {
		return apperrors.TCRERR_MANIFEST_INVALID.Wrap(err)
	}

	existsName, err := u.repo.ExistsName(metadata.Name)
	if err != nil {
		return apperrors.TCRERR_PERSISTER_ERROR.Wrap(err)
	}
	if !existsName {
		return apperrors.TCRERR_NAME_NOT_FOUND
	}

	encodedManifest := base64.StdEncoding.EncodeToString(manifest)
	err = u.repo.PutManifest(metadata, encodedManifest)
	if err != nil {
		return apperrors.TCRERR_PERSISTER_ERROR.Wrap(err)
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
