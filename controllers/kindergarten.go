package controllers

import (
	"encoding/json"
	"errors"
	"kindergarten/models"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	log "github.com/astaxie/beego/logs"
)

type KgController struct {
	beego.Controller
}

// @Title 添加园区
// @Description 添加园区
// @Success 200 {object} models.Kindergarten
// @Param   token header     string    true       "会话token"
// @Param   userid header     int    true       "会话userid"
// @Param   body     body   models.KgAddReq true       "园区名称"
// @Failure 409 园区名已被使用
// @Failure 500 内部错误
// @router / [post]
func (this *KgController) Add() {
	err := checkUserToken(this.Ctx.Input)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}
	req := models.KgAddReq{}
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
	kg, err := models.GetKgByName(req.KgName)
	if err != nil {
		log.Error(err)
		if err != models.ErrNotFound {
			outputInternalError(this.Ctx.Output, err)
			return
		}
	} else {
		this.Ctx.Output.SetStatus(http.StatusConflict)
		this.Ctx.Output.Body([]byte("kindergarten name has been used"))
		return
	}
	kg = &models.Kindergarten{}
	kg.KgName = req.KgName
	id, err := models.InsertKindergarten(kg)

	kg.KgID = id
	this.Ctx.Output.JSON(&kg, false, false)
	return
}

// @Title 修改园区
// @Description 修改园区
// @Success 200 {object} models.KgRes
// @Param   token header     string    true       "会话token"
// @Param   userid header     int    true       "会话userid"
// @Param   body     body   models.KgAddReq true       "园区名称"
// @Failure 409 园区名已被使用
// @Failure 500 内部错误
// @router / [put]
func (this *KgController) Modify() {
	err := checkUserToken(this.Ctx.Input)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}
	req := models.KgModReq{}
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
	kg, err := models.GetKgByName(req.KgName)
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
	kg, err = models.GetKgByID(req.KgID)
	if err != nil {
		if err == models.ErrNotFound {
			outputNotFound(this.Ctx.Output, err)
		} else {
			outputInternalError(this.Ctx.Output, err)
		}
		return
	}
	kg.KgName = req.KgName
	err = models.UpdateKg(kg)
	if err != nil {
		outputInternalError(this.Ctx.Output, err)
		return
	}
	//	res := models.KgRes{}
	//	res.KgID = req.KgID
	//	res.KgName = req.KgName
	this.Ctx.Output.JSON(&kg, false, false)
	return
}

// @Title 查询园区
// @Description 以园区名为条件查询园区
// @Success 200 {object} models.Kindergarten
// @Param   token header     string    true       "会话token"
// @Param   userid header     int    true       "会话userid"
// @Param   name path     string    true       "园区名称"
// @Failure 404 园区不存在
// @Failure 500 内部错误
// @router /kgname/:name [get]
func (this *KgController) Query() {
	err := checkUserToken(this.Ctx.Input)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}

	name := this.Ctx.Input.Param(":name")
	if name == "" {
		outputBadReq(this.Ctx.Output, errors.New("kindergarten name is nil"))
		return
	}
	log.Debug("name", name)
	kg, err := models.GetKgByName(name)
	if err != nil {
		if err != models.ErrNotFound {
			outputInternalError(this.Ctx.Output, err)
			return
		} else {
			log.Debug("kindergarten not found")
			this.Ctx.Output.SetStatus(http.StatusNotFound)
			this.Ctx.Output.Body([]byte("kindergarten not found"))
			return
		}
	}

	this.Ctx.Output.JSON(&kg, false, false)
	return
}

// @Title 查询园区的所有区域
// @Description 查询园区的所有区域
// @Success 200 {object} models.BlocksRes
// @Param   token header     string    true       "会话token"
// @Param   userid header     int    true       "会话userid"
// @Param   kgid path     int    true       "园区ID"
// @Failure 404 园区不存在
// @Failure 500 内部错误
// @router /:kgid/blocks [get]
func (this *KgController) KgBlocks() {
	err := checkUserToken(this.Ctx.Input)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}

	kgid, err := strconv.ParseInt(this.Ctx.Input.Param(":kgid"), 0, 64)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}
	//	log.Debug("name", name)
	_, err = models.GetKgByID(kgid)
	if err != nil {
		if err != models.ErrNotFound {
			outputInternalError(this.Ctx.Output, err)
			return
		} else {
			log.Debug("kindergarten not found")
			this.Ctx.Output.SetStatus(http.StatusNotFound)
			this.Ctx.Output.Body([]byte("kindergarten not found"))
			return
		}
	}
	bs, err := models.GetBlocksByKgID(kgid)
	res := models.BlocksRes{
		Num:    len(bs),
		Blocks: bs,
	}
	this.Ctx.Output.JSON(&res, false, false)
	return
}

// @Title 查询园区的所有设备
// @Description 查询园区的所有区域
// @Success 200 {object} models.DevsRes
// @Param   token header     string    true       "会话token"
// @Param   userid header     int    true       "会话userid"
// @Param   kgid path     int    true       "园区ID"
// @Param   status query     int    false       "设备状态:-2 离线，-1 关，1 开,不上送返回所有状态"
// @Failure 404 园区不存在
// @Failure 500 内部错误
// @router /:kgid/devices [get]
func (this *KgController) KgDevices() {
	err := checkUserToken(this.Ctx.Input)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}

	kgid, err := strconv.ParseInt(this.Ctx.Input.Param(":kgid"), 0, 64)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}
	//	log.Debug("name", name)
	_, err = models.GetKgByID(kgid)
	if err != nil {
		if err != models.ErrNotFound {
			outputInternalError(this.Ctx.Output, err)
			return
		} else {
			log.Debug("kindergarten not found")
			this.Ctx.Output.SetStatus(http.StatusNotFound)
			this.Ctx.Output.Body([]byte("kindergarten not found"))
			return
		}
	}
	status, err := this.GetInt("status", 0)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}
	log.Debug("status=", status)
	ds, err := models.GetDevsByKgID(kgid, status)
	if err != nil {
		log.Error(err)
		outputInternalError(this.Ctx.Output, err)
		return
	}
	res := models.DevsRes{
		Num:     len(ds),
		Devices: ds,
	}
	this.Ctx.Output.JSON(&res, false, false)
	return
}

// @Title 删除园区
// @Description 删除园区
// @Success 200
// @Param   token header     string    true       "会话token"
// @Param   userid header     int    true       "会话userid"
// @Param   kgid     path   int true       "园区ID"
// @Failure 404 园区不存在
// @Failure 500 内部错误
// @router /:kgid [delete]
func (this *KgController) Delete() {
	err := checkUserToken(this.Ctx.Input)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}

	kgid, err := strconv.ParseUint(this.Ctx.Input.Param(":kgid"), 0, 64)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}

	//	err = req.Check()
	//	if err != nil {
	//		log.Error(err)
	//		outputBadReq(this.Ctx.Output, err)
	//		return
	//	}

	err = models.DelKg(int64(kgid))
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

// @Title 获取所有园区
// @Description 获取所有园区
// @Success 200 {object} models.KgsRes
// @Param   token header     string    true       "会话token"
// @Param   userid header     int    true       "会话userid"
// @Failure 500 内部错误
// @router /all [get]
func (this *KgController) All() {
	err := checkUserToken(this.Ctx.Input)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}

	kgs, err := models.GetAllKgs()
	if err != nil {
		outputInternalError(this.Ctx.Output, err)
		return
	}
	kgsRes := models.KgsRes{
		Num: int64(len(kgs)),
		Kgs: kgs,
	}
	this.Ctx.Output.JSON(&kgsRes, false, false)
	return
}
