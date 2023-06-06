package controller

import (
	"github.com/gin-gonic/gin"

	commonctl "github.com/opensourceways/defect-manager/common/controller"
	"github.com/opensourceways/defect-manager/defect/app"
)

type DefectController struct {
	service app.DefectService
}

func AddRouteForDefectController(r *gin.RouterGroup, s app.DefectService) {
	ctl := DefectController{
		service: s,
	}

	r.POST("/v1/defect", ctl.Add)
	r.GET("/v1/defect", ctl.Collect)
	r.POST("/v1/defect/bulletin", ctl.GenerateBulletin)
}

func (ctl DefectController) Add(ctx *gin.Context) {
	commonctl.SendRespOfPost(ctx, "wawa")
}

func (ctl DefectController) Collect(ctx *gin.Context) {

}

func (ctl DefectController) GenerateBulletin(ctx *gin.Context) {

}
