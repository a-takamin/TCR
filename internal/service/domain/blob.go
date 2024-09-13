package domain

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/a-takamin/tcr/internal/interface/persister"
)

type BlobDomain struct {
}

func NewBlobDomain(repository persister.BlobPersister) *BlobDomain {
	return &BlobDomain{}
}

func (s BlobDomain) ValidateNameSpace(namespace string) error {
	matched, _ := regexp.MatchString(`^[a-z0-9]+([._-][a-z0-9]+)*(/[a-z0-9]+([._-][a-z0-9]+)*)*$`, namespace)
	if matched {
		return nil
	}
	return errors.New("name is invalid")
}

func (s BlobDomain) ValidateDigest(digestLike string) error {
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

func (s BlobDomain) ValidateContentRange(contentRangeLike string) error {
	matched, _ := regexp.MatchString(`^[0-9]+-[0-9]+$`, contentRangeLike)
	if !matched {
		return errors.New("Content-Range format is invalid")
	}
	return nil
}

func (s BlobDomain) GetContentRangeStart(contentRange string) (int64, error) {
	ranges := strings.Split(contentRange, "-")
	if len(ranges) != 2 {
		return 0, errors.New("Content-Range format is invalid")
	}
	i, err := strconv.ParseInt(ranges[0], 10, 64)
	if err != nil {
		return 0, errors.New("Content-Range format is invalid")
	}

	return i, nil
}

func (s BlobDomain) GetContentRangeEnd(contentRange string) (int64, error) {
	ranges := strings.Split(contentRange, "-")
	if len(ranges) != 2 {
		return 0, errors.New("Content-Range format is invalid")
	}
	i, err := strconv.ParseInt(ranges[1], 10, 64)
	if err != nil {
		return 0, errors.New("Content-Range format is invalid")
	}

	return i, nil
}
