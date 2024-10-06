package dto

import "io"

type ExistsBlobInput struct {
	Name   string
	Digest string
}

type FindBlobInput struct {
	Name   string
	Digest string
}

type FindChunkedBlobInput struct {
	Name       string
	Uuid       string
	ChunkSeqNo int
}

type FindBlobOutput struct {
	Blob []byte
}

type SaveBlobInput struct {
	Name   string
	Digest string
	Blob   io.Reader
}

type SaveChunkedBlobInput struct {
	Name       string
	Uuid       string
	ChunkSeqNo int
	Blob       io.Reader
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
