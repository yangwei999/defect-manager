package dp

import (
	"errors"
	"regexp"
)

var systemReg = regexp.MustCompile(`^openEuler-`)

type systemVersion string

type SystemVersion interface {
	Version() string
}

func NewSystemVersion(s string) (SystemVersion, error) {
	if !systemReg.MatchString(s) {
		return nil, errors.New("invalid system version")
	}

	return systemVersion(s), nil
}

func (s systemVersion) Version() string {
	return string(s)
}