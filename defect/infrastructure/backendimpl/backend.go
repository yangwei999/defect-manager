package backendimpl

import "github.com/opensourceways/server-common-lib/utils"

func NewBackend(cfg *Config) backendImpl {
	return backendImpl{
		cli: utils.NewHttpClient(3),
	}
}

type backendImpl struct {
	cli utils.HttpClient
}

func (impl backendImpl) MaxBulletinID() (int, error) {

	return 1008, nil
}

func (impl backendImpl) IsDefectPublished([]string) (map[string]bool, error) {

	return nil, nil
}
