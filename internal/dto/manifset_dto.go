package dto

import "github.com/a-takamin/tcr/internal/model"

type GetTagsResponse struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type GetManifestResponse struct {
	Manifest model.Manifest
	Digest   string
}
