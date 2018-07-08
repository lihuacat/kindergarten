package controllers

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"kindergarten/models"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	log "github.com/astaxie/beego/logs"
)

type UserController struct {
	beego.Controller
}

// @Title 用户注册
// @Description 用户注册
// @Success 200 {object} models.RegRes
// @Param   body     body   models.RegReq true       "手机号和密码"
// @Failure 409 手机号已被使用
// @Failure 500 内部错误
// @router / [post]
func (this *UserController) Register() {
	regBody := models.RegReq{}
	json.Unmarshal(this.Ctx.Input.RequestBody, &regBody)
	err := regBody.Check()
	if err != nil {
		log.Error(err)
		this.Ctx.Output.SetStatus(http.StatusBadRequest)
		this.Ctx.Output.Body([]byte(err.Error()))
		return
	}

	_, err = models.GetUserByCellNum(regBody.CellNum)
	if err != nil {
		log.Error(err)
		if err != models.ErrNotFound {
			this.Ctx.Output.SetStatus(http.StatusInternalServerError)
			this.Ctx.Output.Body([]byte(err.Error()))
			return
		}
	} else {
		this.Ctx.Output.SetStatus(http.StatusConflict)
		this.Ctx.Output.Body([]byte(models.ErrCellUsed.Error()))
		return
	}

	newUser := models.User{}
	newUser.CellNum = regBody.CellNum
	//	base64.StdEncoding.EncodeToString(md5.Sum([]byte(regBody.Password))
	md5 := md5.Sum([]byte(regBody.Password))
	newUser.Passwd = base64.StdEncoding.EncodeToString(md5[0:])
	newUser.UserName = regBody.UserName

	id, err := models.InsertUser(&newUser)
	if err != nil {
		log.Error(err)
		this.Ctx.Output.SetStatus(http.StatusInternalServerError)
		this.Ctx.Output.Body([]byte(err.Error()))
		return
	}

	regRes := models.RegRes{}
	regRes.UserID = id
	regRes.UserName = regBody.UserName
	regRes.CellNum = regBody.CellNum
	regRes.Token = models.AddSession(id)
	this.Ctx.Output.JSON(&regRes, false, false)
	return
}

// @Title 用户修改信息
// @Description 用户修改信息
// @Success 200 {object} models.ModifyUserRes
// @Param   body     body   models.ModifyUserReq true       "新的用户信息"
// @Failure 409 手机号已被使用
// @Failure 404 用户不存在
// @Failure 500 内部错误
// @router / [put]
func (this *UserController) ModifyProfile() {
	err := checkUserToken(this.Ctx.Input)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}
	reqBody := models.ModifyUserReq{}
	json.Unmarshal(this.Ctx.Input.RequestBody, &reqBody)
	err = reqBody.Check()
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}
	UserID, _ := strconv.ParseInt(this.Ctx.Input.Header("userid"), 0, 64)
	user, err := models.GetUserByID(UserID)
	if err != nil {
		if err == models.ErrNotFound {
			outputNotFound(this.Ctx.Output, err)
		} else {
			outputInternalError(this.Ctx.Output, err)
		}
		return
	}

	if len(reqBody.NewCellNum) > 0 {
		user.CellNum = reqBody.NewCellNum
		u, err := models.GetUserByCellNum(reqBody.NewCellNum)
		if err != nil {
			log.Error(err)
			if err != models.ErrNotFound {
				this.Ctx.Output.SetStatus(http.StatusInternalServerError)
				this.Ctx.Output.Body([]byte(err.Error()))
				return
			}
		} else {
			if u.UserID != UserID {
				outputInternalError(this.Ctx.Output, err)
				return
			}
		}
	}

	if len(reqBody.NewUserName) > 0 {
		user.UserName = reqBody.NewUserName
	}
	err = models.UpdateUser(user)
	if err != nil {
		outputInternalError(this.Ctx.Output, err)
		return
	}
	res := models.ModifyUserRes{}
	res.UserID = UserID
	res.CellNum = reqBody.NewCellNum
	res.UserName = reqBody.NewUserName
	return
}

// @Title 用户修改密码
// @Description 用户修改密码
// @Success 200
// @Param   token header     string    true       "会话token"
// @Param   userid header     string    true       "会话userid"
// @Param   body     body   models.ModifyPasswdReq true       "新密码和旧密码"
// @Failure 403 原密码错误
// @Failure 404 用户不存在
// @Failure 500 内部错误
// @router /passwd [put]
func (this *UserController) ModifyPasswd() {
	err := checkUserToken(this.Ctx.Input)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}
	req := models.ModifyPasswdReq{}
	json.Unmarshal(this.Ctx.Input.RequestBody, &req)
	err = req.Check()
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}
	UserID, _ := strconv.ParseInt(this.Ctx.Input.Header("userid"), 0, 64)
	user, err := models.GetUserByID(UserID)
	if err != nil {
		log.Error(err)
		if err == models.ErrNotFound {
			outputNotFound(this.Ctx.Output, err)
		} else {
			outputInternalError(this.Ctx.Output, err)
		}
		return
	}
	if user.Passwd != req.OldPasswd {
		this.Ctx.Output.SetStatus(http.StatusBadRequest)
		this.Ctx.Output.Body([]byte("old password error"))
		return
	}
	user.Passwd = req.NewPasswd
	err = models.UpdateUserPasswd(user)
	if err != nil {
		log.Error(err)
		outputInternalError(this.Ctx.Output, err)
		return
	}

	return
}

func outputInternalError(output *context.BeegoOutput, err error) {
	log.Debug(err.Error())
	output.SetStatus(http.StatusInternalServerError)
	output.Body([]byte(err.Error()))
}

func outputBadReq(output *context.BeegoOutput, err error) {
	output.SetStatus(http.StatusBadRequest)
	output.Body([]byte(err.Error()))
}

func outputNotFound(output *context.BeegoOutput, err error) {
	output.SetStatus(http.StatusNotFound)
	output.Body([]byte(err.Error()))
}

func checkUserToken(input *context.BeegoInput) error {
	userID, err := strconv.ParseInt(input.Header("userid"), 0, 64)
	if err != nil {
		log.Error(err)
		return err
	}

	userToken := models.UserToken{userID, input.Header("token")}
	return userToken.CheckToken()
}

// @Title 获取登录用户所在园区
// @Description 获取登录用户所在园区
// @Success 200 {object} models.KgsRes
// @Param   token header     string    true       "会话token"
// @Param   userid header     int    true       "会话userid"
// @Failure 500 内部错误
// @router /kgs [get]
func (this *UserController) UserKgs() {
	err := checkUserToken(this.Ctx.Input)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}

	UserID, _ := strconv.ParseInt(this.Ctx.Input.Header("userid"), 0, 64)
	_, err = models.GetUserByID(UserID)
	if err != nil {
		log.Error(err)
		if err == models.ErrNotFound {
			outputNotFound(this.Ctx.Output, err)
		} else {
			outputInternalError(this.Ctx.Output, err)
		}
		return
	}
	kgs, err := models.GetKgsByUserID(UserID)
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

// @Title 获取登录用户权限下设备类型、状态统计
// @Description 获取登录用户权限下设备类型、状态统计
// @Success 200 {object} models.DevStatRes
// @Param   token header     string    true       "会话token"
// @Param   userid header     int    true       "会话userid"
// @Failure 500 内部错误
// @router /device/statistic [get]
func (this *UserController) DevStatistic() {
	err := checkUserToken(this.Ctx.Input)
	if err != nil {
		log.Error(err)
		outputBadReq(this.Ctx.Output, err)
		return
	}
	UserID, _ := strconv.ParseInt(this.Ctx.Input.Header("userid"), 0, 64)
	_, err = models.GetUserByID(UserID)
	if err != nil {
		log.Error(err)
		if err == models.ErrNotFound {
			outputNotFound(this.Ctx.Output, err)
		} else {
			outputInternalError(this.Ctx.Output, err)
		}
		return
	}
	kgs, err := models.GetKgsByUserID(UserID)

	devStats := make([]*models.DevStat, 0)
	for _, kg := range kgs {
		log.Debug(kg.KgName)
		devStat := &models.DevStat{
			KgID:   kg.KgID,
			KgName: kg.KgName,
		}
		devStat.LightOn, err = models.KgDevOnNum(kg.KgID, "灯")
		if err != nil {
			outputInternalError(this.Ctx.Output, err)
			return
		}
		devStat.AirConditionerOn, err = models.KgDevOnNum(kg.KgID, "空调")
		if err != nil {
			outputInternalError(this.Ctx.Output, err)
			return
		}

		devStats = append(devStats, devStat)
	}

	res := &models.DevStatRes{
		Num:      len(devStats),
		DevStats: devStats,
	}
	this.Ctx.Output.JSON(&res, false, false)
	return
}
