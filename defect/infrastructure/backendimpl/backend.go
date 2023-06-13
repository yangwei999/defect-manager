package backendimpl

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/opensourceways/server-common-lib/utils"
)

var instance *backendImpl

func Init(cfg *Config) {
	instance = &backendImpl{
		cli: utils.NewHttpClient(3),
		cfg: cfg,
	}
}

func Instance() *backendImpl {
	return instance
}

type backendImpl struct {
	cli utils.HttpClient
	cfg *Config
}

type maxIdResult struct {
	Code   int    `json:"code"`
	Result string `json:"result"`
	Msg    string `json:"msg"`
}

type publishedDefectResult struct {
	Code   int      `json:"code"`
	Result []string `json:"result"`
	Msg    string   `json:"msg"`
}

func (impl backendImpl) MaxBulletinID() (maxId int, err error) {
	url := fmt.Sprintf("%s/cve-security-notice-server/securitynotice/getMaxNoticeId?notice_type=bug",
		impl.cfg.Endpoint,
	)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	r, _, err := impl.cli.Download(request)
	if err != nil {
		return
	}

	var res maxIdResult
	if err = json.Unmarshal(r, &res); err != nil {
		return
	}

	if res.Code != 0 {
		err = errors.New(res.Msg)

		return
	}

	t := strings.Split(res.Result, "-")

	return strconv.Atoi(t[len(t)-1])
}

func (impl backendImpl) PublishedDefects() (pub []string, err error) {
	url := fmt.Sprintf("%s/cve-security-notice-server/securitynotice/getPublishedBugs",
		impl.cfg.Endpoint,
	)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	r, _, err := impl.cli.Download(request)
	if err != nil {
		return
	}

	var res publishedDefectResult
	if err = json.Unmarshal(r, &res); err != nil {
		return
	}

	if res.Code != 0 {
		err = errors.New(res.Msg)

		return
	}

	return res.Result, nil
}
