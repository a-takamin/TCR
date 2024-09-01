package repositoryapi

import (
	"github.com/a-takamin/tcr/model"
)

type ManifestRepositoryApi interface {
	GetManifest(metadata model.ManifestMetadata) (model.Manifest, error)
	PutManifest(metadata model.ManifestMetadata, manifest model.Manifest) error
	DeleteManifest(metadata model.ManifestMetadata) error
}
