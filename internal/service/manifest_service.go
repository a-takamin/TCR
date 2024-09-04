package service

import (
	"github.com/a-takamin/tcr/internal/interface/persister"
	"github.com/a-takamin/tcr/internal/model"
	"github.com/a-takamin/tcr/internal/service/utils"
)

type ManifestService struct {
	repo persister.ManifestPersister
}

func NewManifestService(repository persister.ManifestPersister) *ManifestService {
	return &ManifestService{repo: repository}
}

func (s ManifestService) GetManifest(metadata model.ManifestMetadata) (model.Manifest, error) {
	err := utils.ValidateName(metadata.Name)
	if err != nil {
		return model.Manifest{}, err
	}

	manifest, err := s.repo.GetManifest(metadata)
	if err != nil {
		// 空だったら 404 を返す実装を追加する。空以外のエラーもあるので注意。
		return model.Manifest{}, err
	}

	// if c.Request.Header["Content-Type"][0] != manifest.MediaType {
	// 	return model.Manifest{}, errors.New("Content-Type does not match")
	// }

	return manifest, nil
}

func (s ManifestService) PutManifest(metadata model.ManifestMetadata, manifest model.Manifest) error {
	err := utils.ValidateName(metadata.Name)
	if err != nil {
		return err
	}

	err = s.repo.PutManifest(metadata, manifest)
	if err != nil {
		return err
	}
	return nil
}

func (s ManifestService) DeleteManifest(metadata model.ManifestMetadata) error {
	err := utils.ValidateName(metadata.Name)
	if err != nil {
		return err
	}

	err = s.repo.DeleteManifest(metadata)
	if err != nil {
		return err
	}
	return nil

}
