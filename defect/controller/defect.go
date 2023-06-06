package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

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

func (ctl DefectController) Collect(ctx *gin.Context) {

}

func (ctl DefectController) GenerateBulletin(ctx *gin.Context) {

}
