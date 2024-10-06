package domain

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/a-takamin/tcr/internal/apperrors"
	"github.com/a-takamin/tcr/internal/model"
)

// マニフェストの仕様: https://github.com/opencontainers/image-spec/blob/v1.0.1/manifest.md
//
// 未定義のフィールドが含まれていようと、REQUIRED が存在すれば OK らしい
//
// これは Conformance Test のデータからそのように判断した
func ValidateManifest(metadata model.ManifestMetadata, target []byte) error {
	var manifest model.Manifest

	err := json.Unmarshal(target, &manifest)
	if err != nil {
		return fmt.Errorf("manifest is invalid: %w", err)
	}

	if metadata.ContentType != manifest.MediaType {
		return errors.New("Content-Type is invalid")
	}

	// TODO: 今は雑なのできっちりバリデーションする
	if manifest.SchemaVersion == 0 {
		return errors.New("manifest is invalid")
	}

	if manifest.Config.MediaType == "" {
		return errors.New("manifest is invalid")
	}

	return nil
}

func ValidateDigest(digestLike string) error {
	arr := strings.Split(digestLike, ":")
	// digest MUST be "algorithm:encodedstring"
	if len(arr) != 2 {
		return apperrors.ErrInvalidReference
	}
	str := arr[1]
	// TCR では sha256 以外のアルゴリズムを認めていない
	matched, _ := regexp.MatchString(`^[a-f0-9]{64}$`, str)
	if !matched {
		return apperrors.ErrInvalidReference
	}
	return nil
}

func CalcManifestDigest(manifest []byte) (string, error) {
	p := sha256.Sum256(manifest)
	return fmt.Sprintf("sha256:%s", fmt.Sprintf("%x", p)), nil
}

func IsDigest(digestLike string) bool {
	arr := strings.Split(digestLike, ":")
	var str string
	str = arr[0]
	if len(arr) > 1 {
		str = arr[1]
	}
	matched, _ := regexp.MatchString(`^[a-f0-9]{64}$`, str)
	return matched
}

func CalcBlobDigest(blob model.Blob) (string, error) {
	return "", nil
}

func CalcManifestDigestRefactor(manifest []byte) (string, error) {
	// 改行や空白によってハッシュ計算のずれが起らぬように統一する
	var out bytes.Buffer
	err := json.Indent(&out, manifest, "", "\t")
	if err != nil {
		return "", err
	}
	b := out.Bytes()

	p := sha256.Sum256(b)
	return fmt.Sprintf("sha256:%s", fmt.Sprintf("%x", p)), nil

}
