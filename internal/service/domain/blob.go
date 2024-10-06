package domain

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/a-takamin/tcr/internal/apperrors"
)

func ValidateContentRange(contentRangeLike string) error {
	matched, _ := regexp.MatchString(`^[0-9]+-[0-9]+$`, contentRangeLike)
	if !matched {
		return apperrors.ErrInvalidContentRange
	}
	return nil
}

func GetContentRangeStart(contentRange string) (int64, error) {
	ranges := strings.Split(contentRange, "-")
	if len(ranges) != 2 {
		return 0, apperrors.ErrInvalidContentRange
	}
	i, err := strconv.ParseInt(ranges[0], 10, 64)
	if err != nil {
		return 0, apperrors.ErrInvalidContentRange
	}

	return i, nil
}

func GetContentRangeEnd(contentRange string) (int64, error) {
	ranges := strings.Split(contentRange, "-")
	if len(ranges) != 2 {
		return 0, apperrors.ErrInvalidContentRange
	}
	i, err := strconv.ParseInt(ranges[1], 10, 64)
	if err != nil {
		return 0, apperrors.ErrInvalidContentRange
	}

	return i, nil
}
