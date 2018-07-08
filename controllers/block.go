package controllers

import (
	"encoding/json"
	"errors"
	"kindergarten/models"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	//	"github.com/astaxie/beego/context"
	log "github.com/astaxie/beego/logs"
)

type BlockController struct {
	beego.Controller
}

// @Title 添加区域
// @Description 添加区域
// @Success 200 {object} models.Block
// @Param   token header     string    true       "会话token"
// @Param   userid header     int    true       "会话userid"
// @Param   body     body   models.BolckAddReq true       "区域信息"
// @Failure 500 内部错误
// @router / [post]
func (this *BlockController) Add() {
	err := checkUserToken(this.Ctx.Input)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}
	req := models.BolckAddReq{}
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
	_, err = models.GetKgByID(req.KgID)
	if err != nil {
		log.Error(err)
		if err != models.ErrNotFound {
			outputInternalError(this.Ctx.Output, err)
			return
		}
		outputBadReq(this.Ctx.Output, errors.New("kindergarten ID not found"))
		return
	}
	newBlock := models.Block{}
	newBlock.BlockName = req.BlockName
	newBlock.KgID = req.KgID
	newBlock.RmtCtrlID = req.RmtCtrlID
	id, err := models.InsertBlock(&newBlock)
	if err != nil {
		log.Error(err)
		outputInternalError(this.Ctx.Output, err)
		return
	}
	newBlock.BlockID = id
	this.Ctx.Output.JSON(&newBlock, false, false)
	return
}

// @Title 删除区域
// @Description 删除区域
// @Success 200
// @Param   token header     string    true       "会话token"
// @Param   userid header     int    true       "会话userid"
// @Param   blockid     path   int true       "区域ID"
// @Failure 404 区域不存在
// @Failure 500 内部错误
// @router /:blockid [delete]
func (this *BlockController) Delete() {
	err := checkUserToken(this.Ctx.Input)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}
	blockid, err := strconv.ParseUint(this.Ctx.Input.Param(":blockid"), 0, 64)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}

	err = models.DelBlock(int64(blockid))
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

// @Title 查询的所有设备
// @Description 查询区域的所有区域
// @Success 200 {object} models.DevsRes
// @Param   token header     string    true       "会话token"
// @Param   userid header     int    true       "会话userid"
// @Param   blockid path     int    true       "区域ID"
// @Param   status query     int    false       "设备状态:-2 离线，-1 关，1 开,不上送返回所有状态"
// @Failure 404 区域不存在
// @Failure 500 内部错误
// @router /:blockid/devices [get]
func (this *BlockController) BlockDevices() {
	err := checkUserToken(this.Ctx.Input)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}

	blkid, err := strconv.ParseInt(this.Ctx.Input.Param(":blockid"), 0, 64)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}
	//	log.Debug("name", name)
	_, err = models.GetBlockByID(blkid)
	if err != nil {
		if err != models.ErrNotFound {
			outputInternalError(this.Ctx.Output, err)
			return
		} else {
			log.Debug("kindergarten not found")
			this.Ctx.Output.SetStatus(http.StatusNotFound)
			this.Ctx.Output.Body([]byte("block not found"))
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
	ds, err := models.GetDevsByBlkID(blkid, status)
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
