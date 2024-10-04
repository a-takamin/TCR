package persister

import (
	"github.com/a-takamin/tcr/internal/dto"
	"github.com/a-takamin/tcr/internal/model"
)

type ManifestPersister interface {
	// TODO: manifest は string ではなく構造体にしてレイヤー境界の共通言語感を出したい
	// GetManifest(metadata model.ManifestMetadata) (string, error)
	PutManifest(metadata model.ManifestMetadata, manifest string) error
	GetTags(name string) (dto.GetTagsResponse, error)
	// リファクタ後
	ExistsName(name string) (bool, error)
	ExistsManifestByDigest(metadata model.ManifestMetadata) (bool, error)
	ExistsManifestByTag(metadata model.ManifestMetadata) (bool, error)
	GetManifestByDigest(metadata model.ManifestMetadata) (string, error)
	GetManifestByTag(metadata model.ManifestMetadata) (string, error)
	DeleteManifestByDigest(metadata model.ManifestMetadata) error
	DeleteManifestByTag(metadata model.ManifestMetadata) error
}
