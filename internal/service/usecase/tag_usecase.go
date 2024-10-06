package usecase

import (
	"github.com/a-takamin/tcr/internal/interface/persister"
)

type TagUseCase struct {
	// TODO Manifest も Tag も同じ Metadata なのでMetadataPersister とかの方が正しい説
	repo persister.ManifestPersister
}

func NewTagUseCase(repo persister.ManifestPersister) *TagUseCase {
	return &TagUseCase{
		repo: repo,
	}
}
