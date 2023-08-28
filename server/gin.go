package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/server-common-lib/interrupts"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/opensourceways/defect-manager/docs"

	"github.com/opensourceways/defect-manager/defect/app"
	"github.com/opensourceways/defect-manager/defect/controller"
	"github.com/opensourceways/defect-manager/defect/infrastructure/backendimpl"
	"github.com/opensourceways/defect-manager/defect/infrastructure/bulletinimpl"
	"github.com/opensourceways/defect-manager/defect/infrastructure/obsimpl"
	"github.com/opensourceways/defect-manager/defect/infrastructure/producttreeimpl"
	"github.com/opensourceways/defect-manager/defect/infrastructure/repositoryimpl"
)

func StartWebServer(port int, timeout time.Duration) {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(logRequest())

	setRouter(r)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	defer interrupts.WaitForGracefulShutdown()

	interrupts.ListenAndServe(srv, timeout)
}

func setRouter(engine *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Title = "Software Package"
	docs.SwaggerInfo.Description = "set header: 'PRIVATE-TOKEN=xxx'"

	v1 := engine.Group(docs.SwaggerInfo.BasePath)
	setApiV1(v1)

	engine.UseRawPath = true
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func setApiV1(v1 *gin.RouterGroup) {
	controller.AddRouteForDefectController(
		v1, app.NewDefectService(
			repositoryimpl.Instance(),
			producttreeimpl.Instance(),
			bulletinimpl.Instance(),
			backendimpl.Instance(),
			obsimpl.Instance(),
		),
	)
}

func logRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		endTime := time.Now()

		logrus.Infof(
			"| %d | %d | %s | %s |",
			c.Writer.Status(),
			endTime.Sub(startTime),
			c.Request.Method,
			c.Request.RequestURI,
		)
	}
}
