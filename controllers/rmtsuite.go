package controllers

import (
	"errors"
	"kindergarten/models"

	"github.com/astaxie/beego"
	//	"github.com/astaxie/beego/context"
	log "github.com/astaxie/beego/logs"
)

type RmtSuiteController struct {
	beego.Controller
}

// @Title 查询遥控器套件
// @Description 查询遥控器套件
// @Success 200 {object} models.RmtSuite
// @Param   rmtctrlid     path   string true       "遥控器ID"
// @Failure 400 参数为空
// @Failure 500 内部错误
// @router /:rmtctrlid [get]
func (this *RmtSuiteController) Get() {

	log.Debug("RmtSuiteController.get")

	rmtctrlid := this.Ctx.Input.Param(":rmtctrlid")
	if rmtctrlid == "" {
		log.Error("rmtctrlid is nil")
		outputBadReq(this.Ctx.Output, errors.New("rmtctrlid is nil"))
		return
	}
	log.Debug("rmtctrlid=", rmtctrlid)
	res, err := models.GetRmtSuiteByRmtCtrlID(rmtctrlid)
	if err != nil {
		log.Error(err)
		outputInternalError(this.Ctx.Output, err)
		return
	}

	this.Ctx.Output.JSON(&res, false, false)
	return

}
