package dto

type FindBlobUploadProgressInput struct {
	Uuid string
}

type FindBlobUploadProgressOutput struct {
	Uuid         string
	ByteUploaded int64
	NextChunkNo  int
	Digest       string
}

type SaveBlobUploadProgressInput struct {
	Uuid         string
	ByteUploaded int64
	NextChunkNo  int
	Digest       string
}

type DeleteBlobUploadProgressInput struct {
	Uuid string
}
