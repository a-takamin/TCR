package utils

import (
	"errors"
	"regexp"
)

func ValidateName(name string) error {
	matched, _ := regexp.MatchString(`^[a-z0-9]+([._-][a-z0-9]+)*(/[a-z0-9]+([._-][a-z0-9]+)*)*$`, name)
	if matched {
		return nil
	}
	return errors.New("name is invalid")
}
