package apperrors

import "errors"

var ErrManifestNotFound = errors.New("manifest not found")
var ErrBlobNotFound = errors.New("blob not found")
var ErrChunkIsNotInSequence = errors.New("chunk is not in sequence")
var ErrAllChunksAreAlreadyUploaded = errors.New("all chunks are already uploaded")
