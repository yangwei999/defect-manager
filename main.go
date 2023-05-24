package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/opensourceways/server-common-lib/logrusutil"
	liboptions "github.com/opensourceways/server-common-lib/options"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/defect-manager/common/infrastructure/postgres"
	"github.com/opensourceways/defect-manager/config"
	"github.com/opensourceways/defect-manager/defect/domain"
	"github.com/opensourceways/defect-manager/defect/infrastructure/repositoryimpl"
)

type options struct {
	service     liboptions.ServiceOptions
	enableDebug bool
}

func (o *options) Validate() error {
	return o.service.Validate()
}

func gatherOptions(fs *flag.FlagSet, args ...string) options {
	var o options

	o.service.AddFlags(fs)

	fs.BoolVar(
		&o.enableDebug, "enable_debug", false, "whether to enable debug model.",
	)

	fs.Parse(args)

	return o
}

func main() {
	logrusutil.ComponentInit("defect-manager")

	o := gatherOptions(
		flag.NewFlagSet(os.Args[0], flag.ExitOnError),
		os.Args[1:]...,
	)
	if err := o.Validate(); err != nil {
		logrus.Errorf("Invalid options, err:%s", err.Error())

		return
	}

	if o.enableDebug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("debug enabled.")
	}

	// Config
	cfg, err := config.LoadConfig(o.service.ConfigFile)
	if err != nil {
		logrus.Errorf("load config, err:%s", err.Error())

		return
	}

	if err = postgres.Init(&cfg.Postgres); err != nil {
		logrus.Errorf("init db failed, err:%s", err.Error())

		return
	}

	pg := repositoryimpl.NewDefect(&cfg.Config)

	d := domain.Defect{Kernel: "waggg"}

	err = pg.Add(d)
	if err != nil {
		fmt.Println(err)
	}

}
