package dp

import (
	"errors"
	"net/url"
)

type dpUrl string

type URL interface {
	URL() string
}

func NewURL(s string) (URL, error) {
	if s == "" {
		return nil, errors.New("empty url")
	}

	if _, err := url.ParseRequestURI(s); err != nil {
		return nil, errors.New("invalid url")
	}

	return dpUrl(s), nil
}

func (u dpUrl) URL() string {
	return string(u)
}
