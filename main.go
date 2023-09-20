package main

import (
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gin-gonic/gin"
	kafka "github.com/opensourceways/kafka-lib/agent"
	server2 "github.com/opensourceways/server-common-lib/gin"
	"github.com/opensourceways/server-common-lib/logrusutil"
	liboptions "github.com/opensourceways/server-common-lib/options"
	postgres "github.com/opensourceways/server-common-lib/postgre"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/net/context"

	"github.com/opensourceways/defect-manager/config"
	"github.com/opensourceways/defect-manager/defect/app"
	"github.com/opensourceways/defect-manager/defect/controller"
	"github.com/opensourceways/defect-manager/defect/domain/dp"
	"github.com/opensourceways/defect-manager/defect/infrastructure/backendimpl"
	"github.com/opensourceways/defect-manager/defect/infrastructure/bulletinimpl"
	"github.com/opensourceways/defect-manager/defect/infrastructure/obsimpl"
	"github.com/opensourceways/defect-manager/defect/infrastructure/producttreeimpl"
	"github.com/opensourceways/defect-manager/defect/infrastructure/repositoryimpl"
	"github.com/opensourceways/defect-manager/docs"
	"github.com/opensourceways/defect-manager/issue"
	messageserver "github.com/opensourceways/defect-manager/message-server"
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
	log := logrus.NewEntry(logrus.StandardLogger())

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

	if err = obsimpl.Init(&cfg.Obs); err != nil {
		logrus.Errorf("init obs failed, err:%s", err.Error())

		return
	}

	if err = repositoryimpl.Init(&cfg.Config); err != nil {
		logrus.Errorf("init repository failed, err:%s", err.Error())

		return
	}

	// kafka
	if err = kafka.Init(&cfg.Kafka, log, nil, ""); err != nil {
		logrus.Errorf("init kafka failed, err:%s", err.Error())

		return
	}

	defer kafka.Exit()

	dp.Init(cfg.Issue.MaintainVersion)

	backendimpl.Init(&cfg.Backend)

	bulletinimpl.Init(&cfg.Bulletin)

	producttreeimpl.Init(&cfg.ProductTree)

	run(cfg, o)
}

func run(cfg *config.Config, o options) {
	service := app.NewDefectService(
		repositoryimpl.Instance(),
		producttreeimpl.Instance(),
		bulletinimpl.Instance(),
		backendimpl.Instance(),
		obsimpl.Instance(),
	)

	if err := issue.InitEventHandler(&cfg.Issue, service); err != nil {
		logrus.Errorf("init event handler failed, err:%s", err.Error())

		return
	}

	err := messageserver.Init(&cfg.MessageServer, issue.Instance)
	if err != nil {
		logrus.Errorf("init message server failed, err:%s", err.Error())

		return
	}

	// run http server
	server2.StartWebServer(o.service.Port, o.service.GracePeriod, func(engine *gin.Engine) {
		docs.SwaggerInfo.BasePath = "/api"
		docs.SwaggerInfo.Title = "Software Package"
		docs.SwaggerInfo.Description = "set header: 'PRIVATE-TOKEN=xxx'"

		v1 := engine.Group(docs.SwaggerInfo.BasePath)
		controller.AddRouteForDefectController(
			v1, app.NewDefectService(
				repositoryimpl.Instance(),
				producttreeimpl.Instance(),
				bulletinimpl.Instance(),
				backendimpl.Instance(),
				obsimpl.Instance(),
			),
		)
		engine.UseRawPath = true
		engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	})

	wait()
}

func wait() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	defer wg.Wait()

	called := false
	ctx, done := context.WithCancel(context.Background())

	defer func() {
		if !called {
			called = true
			done()
		}
	}()

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()

		select {
		case <-ctx.Done():
			logrus.Info("receive done. exit normally")
			return

		case <-sig:
			logrus.Info("receive exit signal")
			called = true
			done()
			return
		}
	}(ctx)

	<-ctx.Done()
}
