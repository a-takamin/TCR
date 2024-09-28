package dto

type GetTagsResponse struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type GetManifestResponse struct {
	Manifest string
	Digest   string
}
