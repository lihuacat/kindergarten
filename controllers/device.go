package controllers

import (
	"encoding/json"
	"kindergarten/models"
	"time"
	//	"net/http"
	"kindergarten/mqtt"
	"strconv"

	"github.com/astaxie/beego"
	log "github.com/astaxie/beego/logs"
)

type DeviceController struct {
	beego.Controller
}

// @Title 添加设备类型
// @Description 添加设备类型
// @Success 200 {object} models.Device
// @Param   token header     string    true       "会话token"
// @Param   userid header     int    true       "会话userid"
// @Param   body     body   models.DevAddReq true       "设备类型信息"
// @Failure 500 内部错误
// @router / [post]
func (this *DeviceController) Add() {
	err := checkUserToken(this.Ctx.Input)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}
	req := models.DevAddReq{}
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &req)
	if err != nil {
		log.Error(err)
		outputInternalError(this.Ctx.Output, err)
		return
	}

	err = req.Check()
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}

	newDev := models.Device{}
	newDev.DevName = req.DevName
	newDev.BlockID = req.BlockID
	newDev.DevTypeID = req.DevTypeID

	id, err := models.InsertDev(&newDev)
	if err != nil {
		log.Error(err)
		outputInternalError(this.Ctx.Output, err)
		return
	}
	log.Debug("id", id)
	newDev.DevID = id
	this.Ctx.Output.JSON(&newDev, false, false)
	return
}

// @Title 删除设备类型
// @Description 删除设备类型
// @Success 200
// @Param   token header     string    true       "会话token"
// @Param   userid header     int    true       "会话userid"
// @Param   devid     path   int true       "设备ID"
// @Failure 404 设备不存在
// @Failure 500 内部错误
// @router /:devid [delete]
func (this *DeviceController) Delete() {
	err := checkUserToken(this.Ctx.Input)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}
	devid, err := strconv.ParseUint(this.Ctx.Input.Param("devid"), 0, 64)
	if err != nil {
		log.Error(err)
		outputInternalError(this.Ctx.Output, err)
		return
	}

	err = models.DelDevType(int64(devid))
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

// @Title 批量打开/关闭设备
// @Description 批量打开/关闭设备
// @Success 200 {object} models.Devices
// @Param   token header     string    true       "会话token"
// @Param   userid header     int    true       "会话userid"
// @Param   body     body   models.Devices true       "要打开/关闭的设备，可只上送设备ID"
// @Failure 500 内部错误
// @router /turnonoff [post]
func (this *DeviceController) TurnOnOff() {
	err := checkUserToken(this.Ctx.Input)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}

	req := models.Devices{}
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &req)
	if err != nil {
		log.Info(string(this.Ctx.Input.RequestBody))
		log.Error(err)
		outputInternalError(this.Ctx.Output, err)
		return
	}

	for _, dev := range req.Devs {
		err = mqtt.ControlDevice(dev.DevID, dev.Status)
		if err != nil {
			log.Error(err)
		}
	}
	time.Sleep(time.Duration(len(req.Devs)*1) * time.Second)
	for _, dev := range req.Devs {
		models.DevSetStatus(dev.DevID, dev.Status)
	}
	for i, dev := range req.Devs {
		d, _ := models.GetDevByID(dev.DevID)
		req.Devs[i] = d
	}
	this.Ctx.Output.JSON(&req, false, false)
	return
}
