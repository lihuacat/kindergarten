package controllers

import (
	"encoding/json"
	"kindergarten/mqtt"

	"github.com/astaxie/beego"
	log "github.com/astaxie/beego/logs"
)

type TransmitController struct {
	beego.Controller
}

// @Title 给设备转发命令
// @Description 给设备转发命令
// @Success 200 {object} models.Devices
// @Param   token header     string    true       "会话token"
// @Param   userid header     int    true       "会话userid"
// @Param   body     body   mqtt.TransmitContent true
// @Failure 500 内部错误
// @router / [post]
func (this *TransmitController) Relay() {
	err := checkUserToken(this.Ctx.Input)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}

	req := mqtt.TransmitContent{}
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &req)
	if err != nil {
		log.Info(string(this.Ctx.Input.RequestBody))
		log.Error(err)
		outputInternalError(this.Ctx.Output, err)
		return
	}
	log.Debug(string(this.Ctx.Input.RequestBody))
	log.Debug("req.BlockID:",req.BlockID)
	log.Debug("req.Content:",req.Content)
	mqtt.Transmit(req.BlockID, req.Content)

	return

}
