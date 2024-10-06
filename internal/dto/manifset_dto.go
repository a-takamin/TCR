package dto

type GetTagsResponse struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type GetManifestResponse struct {
	Manifest string
	Digest   string
}

type ExistsManifestInput struct {
	Name      string
	Reference string
}

type FindManifestInput struct {
	Name      string
	Reference string
}

type FindManifestOutput struct {
	Name   string
	Tag    string
	Digest string
	// TODO: いったん []byte のままにしているので適切な形に直す
	Manifest []byte
}

type SaveManifestInput struct {
	Name     string
	Tag      string
	Digest   string
	Manifest []byte
}

type DeleteManifestInput struct {
	Name      string
	Reference string
}
