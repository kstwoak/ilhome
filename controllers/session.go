package controllers

import (
	"github.com/astaxie/beego"
	"ilhome/models"
)

type SessionController struct {
	beego.Controller
}

//封装好的返回结构 变成json 返回给前端
func (this *SessionController)RetData(resp interface{}){
	this.Data["json"] = resp
	this.ServeJSON()
}

//api/v1.0/session get  获取用户信息
func (this *SessionController) GetSessionInfo() {
	beego.Info("==============api/v1.0/Session get succ!!=============")
	//返回给前端的map结构体
	resp := make(map[string]interface{})
	//返回错误给前端
	resp["errno"] = models.RECODE_SESSIONERR
	resp["errmsg"] = models.RecodeText(models.RECODE_SESSIONERR)
	defer this.RetData(resp)
	name_map := make(map[string]interface{})
	//将登陆或者注册给当前session存好的name字段返回给前端
	name := this.GetSession("name")
	if name != nil {
		resp["errno"] = models.RECODE_OK
		resp["errmsg"] = models.RecodeText(models.RECODE_OK)
		name_map["name"] = name.(string)
		resp["data"] = name_map
	}
	return
}

////api/v1.0/session delete  删除登陆信息
func (this *SessionController) DelSessionInfo() {
	beego.Info("==============api/v1.0/Session get succ!!=============")
	//返回给前端的map结构体
	resp := make(map[string]interface{})
	//返回正确
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	defer this.RetData(resp)

	this.DelSession("name")
	this.DelSession("user_id")
	this.DelSession("mobile")


	return
}