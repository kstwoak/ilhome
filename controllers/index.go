package controllers

import (
	"github.com/astaxie/beego"
	"ilhome/models"
)

type HousesIndexController struct {
	beego.Controller
}

//封装好的返回结构 变成json 返回给前端
func (this *HousesIndexController)RetData(resp interface{}){
	this.Data["json"] = resp
	this.ServeJSON()
}
//api/v1.0/areas get
func (this *HousesIndexController) GetHousesIndexInfo() {
	beego.Info("==============api/v1.0/HousesIndex get succ!!=============")

	//返回给前端的map结构体
	resp := make(map[string]interface{})
	defer this.RetData(resp)
	//１从缓存中redis读取数据
	//2 如果redis有之前的json字符串数据那么就直接返回给前端
	//3 如果redis没有之前的json字符串数据 那就从mysql中进行查找
	//创建1个orm句柄
	//o:= orm.NewOrm()
	////存放查询到的结果
	//var areas []models.Area
	////创建查询条件
	//qs  := o.QueryTable("area")
	//num ,err := qs.All(&areas)
	//if err != nil {
	//	//返回错误给前端
	//	resp["errno"] =models.RECODE_DBERR
	//	resp["errmsg"] =models.RecodeText(models.RECODE_DBERR)
	//	//c.Data["json"] = resp
	//	//c.ServeJSON()
	//	return
	//}
	//if num == 0{
	//	resp["errno"] =models.RECODE_NODATA
	//	resp["errmsg"] = models.RecodeText(models.RECODE_NODATA)
	//	//c.Data["json"] = resp
	//	//c.ServeJSON()
	//	return
	//}
	//成功
	resp["errno"] =models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	//将封装好的结构体发送给前端
	//c.Data["json"] = resp
	//c.ServeJSON()
	return
}
