package issue

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/opensourceways/server-common-lib/utils"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"
)

var committerInstance *committerCache

type ResContent struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type ResCommitter struct {
	Data struct {
		Maintainers      []string `json:"maintainers"`
		CommitterDetails []struct {
			GiteeId []string `json:"gitee_id"`
			Repo    string   `json:"repo"`
		} `json:"committerDetails"`
	} `json:"data"`
}

type committerCache struct {
	committersOfRepo map[string][]string
	CacheAt          string
}

func InitCommitterInstance() {
	committerInstance = &committerCache{
		committersOfRepo: make(map[string][]string),
	}
}

func (c *committerCache) isCommitter(repo, user string) bool {
	if len(c.committersOfRepo) == 0 || c.CacheAt != time.Now().Format("20060102") {
		c.initCommitterCache()
	}

	v, ok := c.committersOfRepo[repo]
	if !ok {
		return false
	}

	set := sets.NewString(v...)

	return set.Has(user)
}

func (c *committerCache) initCommitterCache() {
	cli := utils.NewHttpClient(3)
	for _, sig := range c.getSig() {
		// Accessing too often can cause 503 errors
		time.Sleep(time.Millisecond * 200)

		url := fmt.Sprintf("https://www.openeuler.org/api-dsapi/query/sig/repo/committers?community=openeuler&sig=%s", sig)

		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			logrus.Errorf("new request of sig %s err: %s", sig, err.Error())
			continue
		}
		r, _, err := cli.Download(request)
		if err != nil {
			logrus.Errorf("get assigner of sig %s err: %s", sig, err.Error())
			continue
		}

		var res ResCommitter
		if err = json.Unmarshal(r, &res); err != nil {
			logrus.Errorf("unmarshal of sig %s err: %s", sig, err.Error())
			continue
		}

		for _, v := range res.Data.CommitterDetails {
			if !strings.Contains(v.Repo, "src-openeuler") {
				continue
			}

			split := strings.Split(v.Repo, "/")
			if len(split) < 2 {
				continue
			}

			c.committersOfRepo[split[1]] = append(res.Data.Maintainers, v.GiteeId...)
		}
	}

	c.CacheAt = time.Now().Format("20060102")
}

func (c *committerCache) getSig() []string {
	url := "https://gitee.com/api/v5/repos/openeuler/community/contents/sig"
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logrus.Errorf("new request of sig url error: %s ", err.Error())

		return nil
	}

	cli := utils.NewHttpClient(3)
	var res []ResContent
	r, _, err := cli.Download(request)
	if err != nil {
		logrus.Errorf("get sig of openeuler error: %s", err.Error())

		return nil
	}

	if err = json.Unmarshal(r, &res); err != nil {
		logrus.Errorf("unmarshal sig error: %s", err.Error())

		return nil
	}

	var sig []string
	for _, v := range res {
		if v.Type == "dir" {
			sig = append(sig, v.Name)
		}
	}

	return sig
}
