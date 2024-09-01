package repository

import "github.com/a-takamin/tcr/model"

type DynamoDBManifestRepository struct {
}

func NewDynamoDBManifestRepository() *DynamoDBManifestRepository {
	return &DynamoDBManifestRepository{}
}

func (r DynamoDBManifestRepository) GetManifest(metadata model.ManifestMetadata) (model.Manifest, error) {
	return model.Manifest{}, nil
}

func (r DynamoDBManifestRepository) PutManifest(metadata model.ManifestMetadata, manifest model.Manifest) error {
	return nil
}
func (r DynamoDBManifestRepository) DeleteManifest(metadata model.ManifestMetadata) error {
	return nil
}
