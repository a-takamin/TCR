package dto

import "io"

type GetBlobInput struct {
	Name   string
	Digest string
}

type DeleteBlobInput struct {
	Name   string
	Digest string
}

type UploadMonolithicBlobInput struct {
	Name          string
	Uuid          string
	Digest        string
	ContentLength int64
	ContentType   string
	Blob          io.ReadCloser
}

type UploadChunkedBlobInput struct {
	Name          string
	Uuid          string
	ContentLength int64
	ContentRange  string
	ContentType   string
	Key           string
	Digest        string
	IsLast        bool
	Blob          io.ReadCloser
}

type BlobUploadProgress struct {
	Uuid         string `json:"Uuid"`
	ByteUploaded int64  `json:"ByteUploaded"`
	NextChunkNo  int    `json:"NextChunkNo"`
	Done         bool   `json:"Done"`
	Digest       string `json:"Digest"`
}

type BlobConcatenateProgress struct {
	Digest string `json:"Digest"`
	Status string `json:"Status"`
}
