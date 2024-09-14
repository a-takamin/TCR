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
	Uuid         string `json:"uuid"`
	ByteUploaded int64  `json:"byte_uploaded"`
	NextChunkNo  int    `json:"next_chunk_no"`
	Done         bool   `json:"done"`
	Digest       string `json:"digest"`
}

type BlobConcatenateProgress struct {
	Digest string `json:"digest"`
	Status string `json:"status"`
}
