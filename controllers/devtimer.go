package controllers

import (
	"encoding/json"
	"kindergarten/models"
	"strconv"

	"github.com/astaxie/beego"
	log "github.com/astaxie/beego/logs"
)

type DevTimerController struct {
	beego.Controller
}

// @Title 添加定时
// @Description 添加定时
// @Param   token header     string    true       "会话token"
// @Param   userid header     int    true       "会话userid"
// @Param   body     body   models.DevTimer true       "设备定时"
// @Success 200 {object} models.DevTimer
// @Failure 500 内部错误
// @router / [post]
func (this *DevTimerController) Add() {

	err := checkUserToken(this.Ctx.Input)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}
	req := &models.AddDevTimerReq{}
	res := &models.DevTimer{}
	log.Debug("this.Ctx.Input.RequestBody:", string(this.Ctx.Input.RequestBody))
	err = json.Unmarshal(this.Ctx.Input.RequestBody, req)
	if err != nil {
		log.Error(err)
		outputInternalError(this.Ctx.Output, err)
		return
	}
	res.DevID = req.DevID
	res.OnOff = req.OnOff
	res.Port = req.Port
	res.TurnHour = req.TurnHour
	res.TurnMinute = req.TurnMinute
	res.TimerID, err = models.AddDevTimer(res)

	if err != nil {
		log.Error(err)
		outputInternalError(this.Ctx.Output, err)
		return
	}

	this.Ctx.Output.JSON(res, false, false)
	return
}

// @Title 删除区域
// @Description 删除区域
// @Success 200
// @Param   token header     string    true       "会话token"
// @Param   userid header     int    true       "会话userid"
// @Param   devtimerid     path   int true       "定时开关ID"
// @Failure 404 区域不存在
// @Failure 500 内部错误
// @router /:devtimerid [delete]
func (this *DevTimerController) Delete() {
	err := checkUserToken(this.Ctx.Input)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}
	devtimerid, err := strconv.ParseUint(this.Ctx.Input.Param(":devtimerid"), 0, 64)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}

	err = models.DelDevTimer(int64(devtimerid))
	if err != nil {
		log.Error(err)
		if err == models.ErrNotFound {
			outputNotFound(this.Ctx.Output, err)
		} else {
			outputInternalError(this.Ctx.Output, err)
		}
		return
	}

	return
}
