package controllers

import (
	"encoding/json"
	"kindergarten/models"

	"github.com/astaxie/beego"
	log "github.com/astaxie/beego/logs"
)

type BlockUserController struct {
	beego.Controller
}

// @Title 添加区域
// @Description 添加区域
// @Success 200 {object} models.BlockUser
// @Param   token header     string    true       "会话token"
// @Param   userid header     int    true       "会话userid"
// @Param   body     body   models.BlockUser true       "区域用户关系"
// @Failure 500 内部错误
// @router / [post]
func (this *BlockUserController) Add() {
	err := checkUserToken(this.Ctx.Input)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}
	req := models.BlockUser{}
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

	_, err = models.GetUserByID(req.UserID)
	if err != nil {
		if err != models.ErrNotFound {
			outputInternalError(this.Ctx.Output, err)
			return
		}
		outputBadReq(this.Ctx.Output, err)
		return
	}
	_, err = models.GetBlockByID(req.BlockID)
	if err != nil {
		if err != models.ErrNotFound {
			outputInternalError(this.Ctx.Output, err)
			return
		}
		outputBadReq(this.Ctx.Output, err)
		return
	}

	err = models.InserBlockUser(&req)
	if err != nil {
		log.Error(err)
		outputInternalError(this.Ctx.Output, err)
		return
	}

	this.Ctx.Output.JSON(&req, false, false)
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
//func (this *UserBlockController) Delete() {
//	err := checkUserToken(this.Ctx.Input)
//	if err != nil {
//		log.Error(err)
//		outputBadReq(this.Ctx.Output, err)
//		return
//	}
//	blockid, err := strconv.ParseUint(this.Ctx.Input.Param("blockid"), 0, 64)
//	if err != nil {
//		log.Error(err)
//		outputBadReq(this.Ctx.Output, err)
//		return
//	}

//	err = models.DelBlock(int64(blockid))
//	if err != nil {
//		log.Error(err)
//		if err == models.ErrNotFound {
//			outputNotFound(this.Ctx.Output, err)
//		} else {
//			outputInternalError(this.Ctx.Output, err)
//		}
//		return
//	}

//	return
//}
