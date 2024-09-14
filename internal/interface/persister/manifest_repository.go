package persister

import (
	"github.com/a-takamin/tcr/internal/dto"
	"github.com/a-takamin/tcr/internal/model"
)

type ManifestPersister interface {
	GetManifest(metadata model.ManifestMetadata) (model.Manifest, error)
	PutManifest(metadata model.ManifestMetadata, manifest model.Manifest) error
	DeleteManifest(metadata model.ManifestMetadata) error
	GetTags(name string) (dto.GetTagsResponse, error)
}
