package utils

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
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
	log.Println(arr)
	var str string
	str = arr[0]
	if len(arr) > 1 {
		str = arr[1]
	}
	log.Println(str)
	matched, _ := regexp.MatchString(`^[a-f0-9]{64}$`, str)
	log.Println("isDigest: ", matched)
	return matched
}
