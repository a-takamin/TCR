package utils

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/a-takamin/tcr/internal/model"
)

func CalcManifestDigest(manifest model.Manifest) (string, error) {
	b, err := json.Marshal(manifest)
	if err != nil {
		return "", err
	}
	p := sha256.Sum256(b)
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
