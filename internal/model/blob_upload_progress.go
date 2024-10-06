package model

// json タグ要る？
type BlobUploadProgress struct {
	Uuid         string `json:"Uuid"`
	ByteUploaded int64  `json:"ByteUploaded"`
	NextChunkNo  int    `json:"NextChunkNo"`
	Done         bool   `json:"Done"`
	Digest       string `json:"Digest"`
}
