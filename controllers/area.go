package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"ilhome/models"
	"encoding/json"
	"time"
)

type AreaController struct {
	beego.Controller
}

//封装好的返回结构 变成json 返回给前端
func (this *AreaController)RetData(resp interface{}){
	this.Data["json"] = resp
	this.ServeJSON()
}

//api/v1.0/areas get
func (this *AreaController) GetAreaInfo() {
	beego.Info("=============api/v1.0/areas get succ!!=============")

	//返回给前端的map结构体
	resp := make(map[string]interface{})
	resp["errno"] =models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	defer this.RetData(resp)

	//0 连接redis数据库
	cache_conn, err := cache.NewCache("redis", `{"key":"ilhome","conn":"127.0.0.1:6379","dbNum":"0"} `)
	if err !=nil {
		beego.Info("cache redis conn err ,err = ",err )
		resp["errno"] =models.RECODE_DBERR
		resp["errmsg"] =models.RecodeText(models.RECODE_DBERR)
		return
	}

	//１从缓存中redis读取数据

	areas_info_value := cache_conn.Get("area_info")
	if areas_info_value !=nil{
		//2 如果redis有之前的json字符串数据那么就直接返回给前端
		//说明area_info key是存在的value就是要返回给前端的json值
		beego.Info("==============get area_info from cache redis succ!!=============")

		var area_info interface{}

		json.Unmarshal(areas_info_value.([]byte), &area_info)
		resp["data"] = area_info
		return
	}

	//3 如果redis没有之前的json字符串数据 那就从mysql中进行查找

	//创建1个orm句柄
	o:= orm.NewOrm()
	//存放查询到的结果
	var areas []models.Area
	//创建查询条件
	qs  := o.QueryTable("area")
	num ,err := qs.All(&areas)
	if err != nil {
		//返回错误给前端
		resp["errno"] =models.RECODE_DBERR
		resp["errmsg"] =models.RecodeText(models.RECODE_DBERR)
		//c.Data["json"] = resp
		//c.ServeJSON()
		return
	}
	if num == 0{
		resp["errno"] =models.RECODE_NODATA
		resp["errmsg"] = models.RecodeText(models.RECODE_NODATA)
		//c.Data["json"] = resp
		//c.ServeJSON()
		return
	}
	//成功

	resp["data"] = areas



	//将areas json字符串 存到 area_info redis 的key中
	areas_info_str ,_ := json.Marshal(areas)
	err = cache_conn.Put("area_info",areas_info_str ,time.Second*3600)
	if err!= nil {
		beego.Info("set area_info --> redis  err =",err)
		resp["errno"] =models.RECODE_DBERR
		resp["errmsg"] =models.RecodeText(models.RECODE_DBERR)
		return
	}

	//将封装好的结构体发送给前端
	//c.Data["json"] = resp
	//c.ServeJSON()
	return
}
