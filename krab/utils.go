package krab

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

var (
	// ErrNoSuchKey is returned if a key of the given name is not found in the store
	ErrNoSuchKey = fmt.Errorf("no key by the given name was found")
	// ErrKeyExists is returned when writing a key would overwrite an existing key
	ErrKeyExists = fmt.Errorf("key by that name already exists, refusing to overwrite")
	// ErrKeyFmt is returned when the key's format is invalid
	ErrKeyFmt = fmt.Errorf("key has invalid format")
)

func validateName(name string) error {
	if name == "" {
		return errors.Wrap(ErrKeyFmt, "key names must be at least one character")
	}

	if strings.Contains(name, "/") {
		return errors.Wrap(ErrKeyFmt, "key names may not contain slashes")
	}

	if strings.HasPrefix(name, ".") {
		return errors.Wrap(ErrKeyFmt, "key names may not begin with a period")
	}

	return nil
}
