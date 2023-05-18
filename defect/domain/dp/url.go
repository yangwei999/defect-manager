package dp

import (
	"errors"
	"regexp"
)

var urlReg = regexp.MustCompile(`^(http|https)://.+`)

type url string

type URL interface {
	URL() string
}

func NewUrl(s string) (URL, error) {
	if !urlReg.MatchString(s) {
		return nil, errors.New("invalid url")
	}

	return url(s), nil
}

func (u url) URL() string {
	return string(u)
}
