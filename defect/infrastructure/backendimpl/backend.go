package backendimpl

import "github.com/opensourceways/server-common-lib/utils"

var instance *backendImpl

func Init(cfg *Config) {
	instance = &backendImpl{
		cli: utils.NewHttpClient(3),
	}
}

func Instance() *backendImpl {
	return instance
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
