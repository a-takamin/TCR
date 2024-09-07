package model

type BlobMetadata struct {
	Name   string
	Digest string
}

type BlobUploadMetadata struct {
	Name          string
	Uuid          string
	Key           string
	Digest        string
	ContentLength int64
	ContentRange  string
	ContentType   string
	// S3 のマルチパートアップロードの URL を発行し、サーバー側でいい感じに保管しておけば
	// ユーザーからのチャンクを受け取ってはその URL にアップロードするだけで簡単にチャンク対応できるが
	// 今回はあえて自前でマルチパートアップロードをやってみる。
	// 具体的には、S3 にチャンクを保存していき、最後のComplete リクエストが来たら Lambda を使って結合する
	IsChunkUpload bool
}

type Blob struct {
}

type BlobUploadProgress struct {
	Name         string `json:"name"`
	ByteUploaded int64  `json:"byte_uploaded"`
	NextChunkNo  int    `json:"next_chunk_no"`
	Done         bool   `json:"done"`
}
