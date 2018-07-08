package controllers

import (
	"encoding/json"
	"kindergarten/models"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	//	"github.com/astaxie/beego/context"
	log "github.com/astaxie/beego/logs"
)

type DevTypeController struct {
	beego.Controller
}

// @Title 添加设备类型
// @Description 添加设备类型
// @Success 200 {object} models.DevType
// @Param   token header     string    true       "会话token"
// @Param   userid header     int    true       "会话userid"
// @Param   body     body   models.DevTypeAddReq true       "设备类型信息"
// @Failure 500 内部错误
// @router / [post]
func (this *DevTypeController) Add() {
	err := checkUserToken(this.Ctx.Input)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}
	req := models.DevTypeAddReq{}
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
	_, err = models.GetDevTypeByName(req.DevTypeName)
	if err != nil {
		if err != models.ErrNotFound {
			outputInternalError(this.Ctx.Output, err)
			return
		}
	} else {
		this.Ctx.Output.SetStatus(http.StatusConflict)
		this.Ctx.Output.Body([]byte("kindergarten name has been used"))
		return
	}
	newType := models.DevType{}
	newType.DevTypeName = req.DevTypeName

	id, err := models.InsertDevType(&newType)

	newType.DevTypeID = id
	this.Ctx.Output.JSON(&newType, false, false)
	return
}

// @Title 查询设备类型
// @Description 查询所有设备类型
// @Success 200 {object} models.DevTypesRes
// @Param   token header     string    true       "会话token"
// @Param   userid header     int    true       "会话userid"
// @Failure 500 内部错误
// @router /all [get]
func (this *DevTypeController) All() {
	err := checkUserToken(this.Ctx.Input)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}

	devTypes, err := models.GetDevTypes()
	if err != nil {
		this.Ctx.Output.SetStatus(http.StatusInternalServerError)
		this.Ctx.Output.Body([]byte(err.Error()))
		return
	}
	res := models.DevTypesRes{
		Num:      len(devTypes),
		DevTypes: devTypes,
	}
	this.Ctx.Output.JSON(&res, false, false)
	return
}

// @Title 删除设备类型
// @Description 删除设备类型
// @Success 200
// @Param   token header     string    true       "会话token"
// @Param   userid header     int    true       "会话userid"
// @Param   devtypeid     path   int true       "设备类型ID"
// @Failure 404 设备类型不存在
// @Failure 500 内部错误
// @router /:devtypeid [delete]
func (this *DevTypeController) Delete() {
	err := checkUserToken(this.Ctx.Input)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}
	devtypeid, err := strconv.ParseUint(this.Ctx.Input.Param("devtypeid"), 0, 64)
	if err != nil {
		log.Error(err)
		outputInternalError(this.Ctx.Output, err)
		return
	}

	err = models.DelDevType(int64(devtypeid))
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
