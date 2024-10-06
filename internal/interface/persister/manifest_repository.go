package persister

import (
	"github.com/a-takamin/tcr/internal/dto"
)

type ManifestPersister interface {
	// リファクタ
	// 次の段階: なるべきエンティティ（ドメインオブジェクト）を引数や戻り値で扱うようにする
	// つまり、ユースケースにドメインオブジェクトをそのまま永続化しているように感じさせる
	ExistsManifest(input dto.ExistsManifestInput) (bool, error)
	FindManifest(input dto.FindManifestInput) (dto.FindManifestOutput, error)
	SaveManifest(input dto.SaveManifestInput) error
	DeleteManifest(input dto.DeleteManifestInput) error
	// tag
	GetTags(name string) (dto.GetTagsResponse, error)
}
