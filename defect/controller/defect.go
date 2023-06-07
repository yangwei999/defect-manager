package controller

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"

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

	r.POST("/v1/defect", ctl.Save)
	r.GET("/v1/defect", ctl.Collect)
	r.POST("/v1/defect/bulletin", ctl.GenerateBulletin)
}

// Save
// @Summary add or update a defect
// @Description add or update a defect
// @Tags  Defect
// @Accept json
// @Param	param  body	 defectRequest	 true	"body of a defect"
// @Success 201 {object} ResponseData
// @Failure 400 {object} ResponseData
// @Router /v1/defect [post]
func (ctl DefectController) Save(ctx *gin.Context) {
	var req defectRequest
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, err)

		return
	}

	if err := ctl.service.SaveDefects(cmd); err != nil {
		commonctl.SendFailedResp(ctx, "", err)
	} else {
		commonctl.SendRespOfPost(ctx, "")
	}
}

// Collect
// @Summary collect information of some defects
// @Description collect information of some defects
// @Tags  Defect
// @Accept json
// @Param	date  query string	 true	"collect defects after the date"
// @Success 200 {object} []app.CollectDefectsDTO
// @Failure 400 {object} ResponseData
// @Router /v1/defect [get]
func (ctl DefectController) Collect(ctx *gin.Context) {
	date, err := time.Parse("2006-01-02", ctx.Query("date"))
	if err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	if v, err := ctl.service.CollectDefects(date); err != nil {
		commonctl.SendFailedResp(ctx, "", err)
	} else {
		commonctl.SendRespOfGet(ctx, v)
	}
}

// GenerateBulletin
// @Summary generate security bulletin for some defects
// @Description generate security bulletin for some defects
// @Tags  Defect
// @Accept json
// @Param	param  body	 defectRequest	 true	"body of some issue number"
// @Success 201 {object} ResponseData
// @Failure 400 {object} ResponseData
// @Router /v1/defect/bulletin [post]
func (ctl DefectController) GenerateBulletin(ctx *gin.Context) {
	var req bulletinRequest
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		commonctl.SendBadRequestBody(ctx, err)

		return
	}

	commonctl.SendRespOfPost(ctx, "Processing: Data is being prepared, please wait patiently\n")

	go func() {
		logrus.Infof("generate bulletin processing of %v", req.IssueNumber)

		if err := ctl.service.GenerateBulletins(req.IssueNumber); err != nil {
			logrus.Errorf("generate bulletin of %v err: %s", req.IssueNumber, err.Error())
		} else {
			logrus.Infof("generate bulletin success of %v", req.IssueNumber)
		}
	}()

}
