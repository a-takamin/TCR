package persister

import "github.com/a-takamin/tcr/internal/dto"

type RepositoryPersister interface {
	ExistsRepository(input dto.ExistsRepositoryInput) (bool, error)
	SaveRepository(input dto.SaveRepositoryInput) error
	DeleteRepository(input dto.DeleteRepositoryInput) error
}
