package usecase

import (
	"github.com/a-takamin/tcr/internal/dto"
	"github.com/a-takamin/tcr/internal/interface/persister"
	"github.com/a-takamin/tcr/internal/service/domain"
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

func (u TagUseCase) GetTags(name string) (dto.GetTagsResponse, error) {
	err := domain.ValidateNameSpace(name)
	if err != nil {
		return dto.GetTagsResponse{}, err
	}
	return u.repo.GetTags(name)
}
