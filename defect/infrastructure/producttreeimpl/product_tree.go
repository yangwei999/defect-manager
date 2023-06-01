package producttreeimpl

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/opensourceways/robot-gitee-lib/client"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/defect-manager/defect/domain"
	"github.com/opensourceways/defect-manager/defect/domain/dp"
)

func NewProductTree(cfg *Config) *ProductTreeImpl {
	return &ProductTreeImpl{
		cli: client.NewClient(func() []byte {
			return []byte(cfg.Token)
		}),
		cfg:                 cfg,
		rpmCache:            make(map[string][]byte),
		rpmOfComponentCache: make(map[string]string),
	}
}

type ProductTreeImpl struct {
	cli client.Client
	cfg *Config

	rpmCache            map[string][]byte
	rpmOfComponentCache map[string]string
}

func (impl *ProductTreeImpl) CleanCache() {
	impl.rpmCache = make(map[string][]byte)
	impl.rpmOfComponentCache = make(map[string]string)
}

func (impl *ProductTreeImpl) GetTree(component string, versions []dp.SystemVersion) (domain.ProductTree, error) {
	if len(impl.rpmCache) == 0 {
		impl.initRPMCache()
	}

	affectedRPM := make(map[string]string)
	for _, v := range versions {
		key := fmt.Sprintf("%s_%s", component, v.String())
		rpm, ok := impl.rpmOfComponentCache[key]
		if !ok {
			rpm = impl.parseRPM(component, v.String())
			impl.rpmOfComponentCache[key] = rpm
		}

		affectedRPM[v.String()] = rpm
	}

	return impl.buildTree(affectedRPM), nil
}

func (impl *ProductTreeImpl) parseRPM(component, version string) string {
	buf := bytes.NewBuffer(impl.rpmCache[version])

	var rpmOfComponent string
	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			logrus.Errorf("error reading %s from %s error %s", component, version, err.Error())
			continue
		}

		split := strings.Split(line, ",")
		if len(split) != 3 {
			logrus.Errorf("the format of line error")
			continue
		}

		if split[1] == component {
			rpmOfComponent = split[2]
			break
		}
	}

	return rpmOfComponent
}

func (impl *ProductTreeImpl) initRPMCache() {
	for version := range dp.MaintainVersion {
		for {
			content, err := impl.cli.GetPathContent(
				impl.cfg.PkgRPM.Org,
				impl.cfg.PkgRPM.Repo,
				fmt.Sprintf("%s%s.csv", impl.cfg.PkgRPM.PathPrefix, version),
				impl.cfg.PkgRPM.Branch,
			)
			if err != nil {
				logrus.Errorf("get content of %s error %s", version, err.Error())
				time.Sleep(time.Second * 3)
				continue
			}

			decodeContent, err := base64.StdEncoding.DecodeString(content.Content)
			if err != nil {
				logrus.Errorf("base64decode content of %s error %s", version, err.Error())
				time.Sleep(time.Second * 3)
				continue
			}

			impl.rpmCache[version.String()] = decodeContent

			break
		}
	}
}

func (impl *ProductTreeImpl) buildTree(affectedRPM map[string]string) domain.ProductTree {
	tree := make(map[string][]domain.Product)
	for version, rpms := range affectedRPM {

		rpmSlice := strings.Fields(rpms)
		for _, rpm := range rpmSlice {
			t := strings.Split(rpm, ".")
			arch := t[len(t)-2]
			productId := strings.Join(t[:len(t)-3], ".")

			product := domain.Product{
				ID:       productId,
				CPE:      version,
				FullName: rpm,
			}

			tree[arch] = append(tree[arch], product)
		}
	}

	return tree
}
