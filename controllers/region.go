package controllers

import (
	"encoding/json"
	"kindergarten/models"
	"strconv"

	"github.com/astaxie/beego"
	log "github.com/astaxie/beego/logs"
)

type RegionController struct {
	beego.Controller
}

// @Title 添加区域
// @Description 添加区域
// @Success 200 {object} models.Region
// @Param   token header     string    true       "会话token"
// @Param   userid header     int    true       "会话userid"
// @Param   body     body   models.RegionReq true       "区域信息"
// @Failure 500 内部错误
// @router / [post]
func (this *RegionController) Add() {
	err := checkUserToken(this.Ctx.Input)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}

	req := models.RegionReq{}
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

	region := &models.Region{}
	region.KgID = req.KgID
	region.Name = req.Name
	region.ID, err = models.AddRegion(region)

	if err != nil {
		log.Error(err)
		outputInternalError(this.Ctx.Output, err)
		return
	}

	this.Ctx.Output.JSON(region, false, false)

	return

}

// @Title 设置region和block关系
// @Description 设置region和block关系
// @Success 200
// @Param   token header     string    true       "会话token"
// @Param   userid header     int    true       "会话userid"
// @Param   regionid    path   int true       "region ID"
// @Param   blockid     path   int true       "block ID"
// @Failure 500 内部错误
// @router /:regionid/block/:blockid [post]
func (this *RegionController) AddBlock() {
	err := checkUserToken(this.Ctx.Input)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}

	regionid, err := strconv.ParseInt(this.Ctx.Input.Param(":regionid"), 0, 64)
	if err != nil {
		log.Error(err)
		outputInternalError(this.Ctx.Output, err)
		return
	}

	blockid, err := strconv.ParseInt(this.Ctx.Input.Param(":blockid"), 0, 64)
	if err != nil {
		log.Error(err)
		outputInternalError(this.Ctx.Output, err)
		return
	}

	models.AddRegionBlock(regionid, blockid)
	return
}

// @Title 删除region
// @Description 删除region
// @Success 200
// @Param   token header     string    true       "会话token"
// @Param   userid header     int    true       "会话userid"
// @Param   regionid    path   int true       "region ID"
// @Failure 500 内部错误
// @router /:regionid [delete]
func (this *RegionController) Delete() {
	err := checkUserToken(this.Ctx.Input)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}

	regionid, err := strconv.ParseInt(this.Ctx.Input.Param(":regionid"), 0, 64)
	if err != nil {
		log.Error(err)
		outputInternalError(this.Ctx.Output, err)
		return
	}

	err = models.DelRegion(regionid)
	if err != nil {
		log.Error(err)
		outputInternalError(this.Ctx.Output, err)
		return
	}

	return
}

// @Title 查询region的所有block
// @Description 查询region的所有block
// @Success 200 {object} models.BlocksRes
// @Param   token header     string    true       "会话token"
// @Param   userid header     int    true       "会话userid"
// @Param   blockid path     int    true       "区域ID"
// @Failure 404 区域不存在
// @Failure 500 内部错误
// @router /:regionid/blocks [get]
func (this *RegionController) BlockRegions() {

	err := checkUserToken(this.Ctx.Input)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}

	regionid, err := strconv.ParseInt(this.Ctx.Input.Param(":regionid"), 0, 64)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}

	blocks, err := models.GetBlocksByRegionID(regionid)
	if err != nil {
		log.Error(err)
		outputInternalError(this.Ctx.Output, err)
		return
	}
	res := &models.BlocksRes{}
	res.Num = len(blocks)
	res.Blocks = blocks

	this.Ctx.Output.JSON(res, false, false)
	return
}
