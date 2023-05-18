package dp

import "errors"

const (
	aarch64 = "aarch64"
	src     = "src"
	x86_64  = "x86_64"
	noarch  = "noarch"
)

var validArch = map[string]bool{
	aarch64: true,
	src:     true,
	x86_64:  true,
	noarch:  true,
}

type arch string

type Arch interface {
	Arch() string
}

func NewArch(s string) (Arch, error) {
	if !validArch[s] {
		return nil, errors.New("invalid arch")
	}

	return arch(s), nil
}

func (a arch) Arch() string {
	return string(a)
}
