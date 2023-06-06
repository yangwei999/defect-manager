package dp

import (
	"errors"
)

const (
	openeuler2003SP1 = "openEuler-20.03-LTS-SP1"
	openeuler2003SP3 = "openEuler-20.03-LTS-SP3"
	openeuler2203    = "openEuler-22.03-LTS"
	openeuler2203SP1 = "openEuler-22.03-LTS-SP1"
)

var MaintainVersion = map[SystemVersion]bool{
	systemVersion(openeuler2003SP1): true,
	systemVersion(openeuler2003SP3): true,
	systemVersion(openeuler2203):    true,
	systemVersion(openeuler2203SP1): true,
}

type systemVersion string

type SystemVersion interface {
	String() string
}

func NewSystemVersion(s string) (SystemVersion, error) {
	// MaintainVersion is not used for validation because
	// there is an error reading old data from the database when maintainVersion changes
	if s == "" {
		return nil, errors.New("invalid system version")
	}

	return systemVersion(s), nil
}

func (s systemVersion) String() string {
	return string(s)
}
