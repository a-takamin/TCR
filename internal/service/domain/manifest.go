package domain

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/a-takamin/tcr/internal/model"
)

func ValidateNameSpace(namespace string) error {
	matched, _ := regexp.MatchString(`^[a-z0-9]+([._-][a-z0-9]+)*(/[a-z0-9]+([._-][a-z0-9]+)*)*$`, namespace)
	if matched {
		return nil
	}
	return errors.New("name is invalid")
}

func ValidateDigest(digestLike string) error {
	arr := strings.Split(digestLike, ":")
	// digest MUST be "algorithm:encodedstring"
	if len(arr) != 2 {
		return errors.New("digest is invalid")
	}
	str := arr[1]
	// TCR では sha256 以外のアルゴリズムを認めていない
	matched, _ := regexp.MatchString(`^[a-f0-9]{64}$`, str)
	if !matched {
		return errors.New("digest is invalid")
	}
	return nil
}

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
