package dto

type GetTagsResponse struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}
