package apperrors

type OCIErrorResponse struct {
	Errors []OCIError `json:"errors"`
}

type OCIError struct {
	ErrorCode    string `json:"code"`
	ErrorMessage string `json:"message"`
	Detail       string `json:"detail"`
}

func (e OCIError) CreateResponse(detail string) OCIErrorResponse {
	if detail != "" {
		e.Detail = detail
	}
	return OCIErrorResponse{
		Errors: []OCIError{e},
	}
}

var BLOB_UNKNOWN = &OCIError{ErrorCode: "BLOB_UNKNOWN", ErrorMessage: "blob unknown to registry"}
var BLOB_UPLOAD_INVALID = &OCIError{ErrorCode: "BLOB_UPLOAD_INVALID", ErrorMessage: "blob upload invalid"}
var BLOB_UPLOAD_UNKNOWN = &OCIError{ErrorCode: "BLOB_UPLOAD_UNKNOWN", ErrorMessage: "blob upload unknown to registry"}
var DIGEST_INVALID = &OCIError{ErrorCode: "DIGEST_INVALID", ErrorMessage: "provided digest did not match uploaded content"}
var MANIFEST_BLOB_UNKNOWN = &OCIError{ErrorCode: "MANIFEST_BLOB_UNKNOWN", ErrorMessage: "blob unknown to registry"}
var MANIFEST_INVALID = &OCIError{ErrorCode: "MANIFEST_INVALID", ErrorMessage: "manifest invalid"}
var MANIFEST_UNKNOWN = &OCIError{ErrorCode: "MANIFEST_UNKNOWN", ErrorMessage: "manifest unknown"}
var MANIFEST_UNVERIFIED = &OCIError{ErrorCode: "MANIFEST_UNVERIFIED", ErrorMessage: "manifest failed signature verification"}
var NAME_INVALID = &OCIError{ErrorCode: "NAME_INVALID", ErrorMessage: "invalid repository name"}
var NAME_UNKNOWN = &OCIError{ErrorCode: "NAME_UNKNOWN", ErrorMessage: "repository name not known to registry"}
var SIZE_INVALID = &OCIError{ErrorCode: "SIZE_INVALID", ErrorMessage: "provided length did not match content length"}
var TAG_INVALID = &OCIError{ErrorCode: "TAG_INVALID", ErrorMessage: "manifest tag did not match URI"}
var UNAUTHORIZED = &OCIError{ErrorCode: "UNAUTHORIZED", ErrorMessage: "authentication required"}
var DENIED = &OCIError{ErrorCode: "DENIED", ErrorMessage: "requested access to the resource is denied"}
var UNSUPPORTED = &OCIError{ErrorCode: "UNSUPPORTED", ErrorMessage: "The operation is unsupported"}
