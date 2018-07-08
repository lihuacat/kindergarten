package controllers

import (
	"encoding/base64"
	"encoding/json"
	//	"fmt"
	"kindergarten/models"
	"strconv"

	"crypto/md5"
	"net/http"

	"github.com/astaxie/beego"
	log "github.com/astaxie/beego/logs"
)

type SessionController struct {
	beego.Controller
}

// @Title 用户登录
// @Description 用户登录
// @Success 200 {object} models.LoginRes
// @Param   body     body   models.LoginReq true       "手机号和密码"
// @Failure 404 用户不存在
// @Failure 500 内部错误
// @router / [post]
func (this *SessionController) Post() {

	loginBody := models.LoginReq{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &loginBody)
	if err != nil {
		log.Error(err)
		log.Debug(string(this.Ctx.Input.RequestBody))
		this.Ctx.Output.SetStatus(http.StatusInternalServerError)
		this.Ctx.Output.Body([]byte(err.Error()))
		return
	}
	err = loginBody.Check()
	if err != nil {
		this.Ctx.Output.SetStatus(http.StatusBadRequest)
		this.Ctx.Output.Body([]byte(err.Error()))
		return
	}
	u, err := models.GetUserByCellNum(loginBody.CellNum)
	if err != nil {
		if err == models.ErrNotFound {
			this.Ctx.Output.SetStatus(http.StatusNotFound)
			this.Ctx.Output.Body([]byte("user not found"))
		} else {
			this.Ctx.Output.SetStatus(http.StatusInternalServerError)
			this.Ctx.Output.Body([]byte(err.Error()))
		}
		return
	}
	md5 := md5.Sum([]byte(loginBody.Password))
	//	log.Debug(fmt.Sprintf("%x", md5))
	if u.Passwd != base64.StdEncoding.EncodeToString(md5[0:]) {
		this.Ctx.Output.SetStatus(http.StatusBadRequest)
		this.Ctx.Output.Body([]byte("password error"))

		return
	}
	token := models.AddSession(u.UserID)
	res := models.LoginRes{
		UserID:   u.UserID,
		UserName: u.UserName,
		Token:    token,
	}

	this.Ctx.Output.JSON(&res, false, false)
	return
}

// @Title 用户退出登录
// @Description 用户退出登录
// @Success 200
// @Param   token header     string    true       "会话token"
// @Param   userid header     int    true       "会话userid"
// @Failure 403 userid或token为空
// @Failure 500 内部错误
// @router / [delete]
func (this *SessionController) Delete() {
	err := checkUserToken(this.Ctx.Input)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}
	userid, _ := strconv.Atoi(this.Ctx.Input.Header("userid"))

	models.DelSession(int64(userid))
	return
}
