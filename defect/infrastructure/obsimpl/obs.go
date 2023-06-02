package obsimpl

import (
	"bytes"
	"fmt"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"

	"github.com/opensourceways/defect-manager/utils"
)

func NewObs(cfg *Config) (obsImpl, error) {
	cli, err := obs.New(cfg.AccessKey, cfg.SecretKey, cfg.Endpoint)
	if err != nil {
		return obsImpl{}, err
	}

	return obsImpl{
		cfg: cfg,
		cli: cli,
	}, nil
}

type obsImpl struct {
	cfg *Config
	cli *obs.ObsClient
}

func (impl obsImpl) Upload(fileName string, data []byte) error {
	input := &obs.PutObjectInput{}
	input.Bucket = impl.cfg.Bucket
	input.Key = fmt.Sprintf("%s/%s/%s", impl.cfg.Directory, utils.Date(), fileName)
	input.Body = bytes.NewReader(data)

	_, err := impl.cli.PutObject(input)

	return err
}
