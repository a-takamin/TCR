package model

type ManifestMetadata struct {
	Name      string
	Reference string
}

type Manifest struct {
	SchemaVersion int               `json:"schemaVersion"`
	MediaType     string            `json:"mediaType"`
	ArtifactType  string            `json:"artifactType"`
	Config        Descriptor        `json:"config"`
	Layers        []Descriptor      `json:"layers"`
	Subject       Descriptor        `json:"subject"`
	Annotations   map[string]string `json:"annotations"`
}

type Descriptor struct {
	MediaType    string            `json:"mediaType"`
	Digest       string            `json:"digest"`
	Size         int64             `json:"size"`
	Urls         []string          `json:"urls"`
	Annotations  map[string]string `json:"annotations"`
	Data         string            `json:"data"`
	ArtifactType string            `json:"artifactType"`
}
