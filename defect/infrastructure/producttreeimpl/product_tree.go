package producttreeimpl

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/opensourceways/robot-gitee-lib/client"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/defect-manager/defect/domain"
	"github.com/opensourceways/defect-manager/defect/domain/dp"
)

func NewProductTree(cfg *Config) *productTreeImpl {
	return &productTreeImpl{
		cli: client.NewClient(func() []byte {
			return []byte(cfg.Token)
		}),
		cfg:                 cfg,
		rpmCache:            make(map[string][]byte),
		rpmOfComponentCache: make(map[string]string),
	}
}

type productTreeImpl struct {
	cli client.Client
	cfg *Config

	rpmCache            map[string][]byte
	rpmOfComponentCache map[string]string

	lock sync.Mutex
	wg   sync.WaitGroup
}

func (impl *productTreeImpl) CleanCache() {
	impl.rpmCache = make(map[string][]byte)
	impl.rpmOfComponentCache = make(map[string]string)
}

func (impl *productTreeImpl) GetTree(component string, versions []dp.SystemVersion) (domain.ProductTree, error) {
	impl.initRPMCache()

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

func (impl *productTreeImpl) parseRPM(component, version string) string {
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

func (impl *productTreeImpl) initRPMCache() {
	impl.lock.Lock()
	if len(impl.rpmCache) == len(dp.MaintainVersion) {
		return
	}
	defer impl.lock.Unlock()

	for version := range dp.MaintainVersion {
		v := version.String()
		impl.wg.Add(1)
		go func() {
			impl.fetchRPMData(v)
			impl.wg.Done()
		}()
	}

	impl.wg.Wait()
}

func (impl *productTreeImpl) fetchRPMData(version string) {
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

		impl.rpmCache[version] = decodeContent

		break
	}
}

func (impl *productTreeImpl) buildTree(affectedRPM map[string]string) domain.ProductTree {
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
