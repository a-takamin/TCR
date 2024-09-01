package service

import (
	"github.com/a-takamin/tcr/interface/persister"
	"github.com/a-takamin/tcr/model"
)

type ManifestService struct {
	repo persister.ManifestPersister
}

func NewManifestService(repository persister.ManifestPersister) *ManifestService {
	return &ManifestService{repo: repository}
}

func (s ManifestService) GetManifest(metadata model.ManifestMetadata) (model.Manifest, error) {
	// 	name := c.Param("name")
	// 	// reference := c.Param("reference")
	// 	match, _ := regexp.MatchString("[a-z0-9]+([._-][a-z0-9]+)*(/[a-z0-9]+([._-][a-z0-9]+)*)*", name)
	// 	if !match {
	// 		return model.Manifest{}, errors.New("requested manifest does not exist")
	// 	}

	// manifest, err := repositories.GetManifest()
	// if err != nil {
	// 	// 空だったら 404 を返す実装を追加する。空以外のエラーもあるので注意。
	// 	return model.Manifest{}, err
	// }

	// 	if c.Request.Header["Content-Type"][0] != manifest.MediaType {
	// 		return model.Manifest{}, errors.New("Content-Type does not match")
	// 	}

	return model.Manifest{}, nil
}
