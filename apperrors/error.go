package apperrors

import "errors"

var ErrManifestNotFound = errors.New("manifest not found")
var ErrBlobNotFound = errors.New("blob not found")
