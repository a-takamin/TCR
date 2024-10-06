package domain

import (
	"regexp"

	"github.com/a-takamin/tcr/internal/apperrors"
)

func ValidateName(name string) error {
	matched, _ := regexp.MatchString(`^[a-z0-9]+([._-][a-z0-9]+)*(/[a-z0-9]+([._-][a-z0-9]+)*)*$`, name)
	if matched {
		return nil
	}
	return apperrors.ErrInvalidName
}
