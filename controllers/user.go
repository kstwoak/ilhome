package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"ilhome/models"
	"regexp"
	"path"
)

type UserController struct {
	beego.Controller
}

//封装好的返回结构 变成json 返回给前端
func (this *UserController)RetData(resp interface{}){
	this.Data["json"] = resp
	this.ServeJSON()
}

//api/v1.0/users post
//{
//mobile: "123",
//password: "123",
//sms_code: "123"
//}

func (this *UserController) Reg() {
	beego.Info("==============api/v1.0/user post succ!!=============")

	//返回给前端的map结构体
	resp := make(map[string]interface{})

	resp["errno"] =models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	defer this.RetData(resp)

	//1获得浏览器请求json数据信息 post数据
	var regRequestMap = make(map[string]interface{})

	json.Unmarshal(this.Ctx.Input.RequestBody, &regRequestMap) //获取前端发送过来的json字符串并解析


	beego.Info("mobile = ", regRequestMap["mobile"])
	beego.Info("password = ", regRequestMap["password"])
	beego.Info("sms_code = ", regRequestMap["sms_code"])
	//2判断数据的合法性
	if regRequestMap["mobile"] == "" || regRequestMap["password"] == "" || regRequestMap["sms_code"] == "" {
		resp["errno"] = models.RECODE_REQERR
		resp["errmsg"] = models.RecodeText(models.RECODE_REQERR)
		return
	}
	//errmobile := regexp.MatchString("0?(13|14|15|17|18|19)[0-9]{9}", "seafood")
	myreg := regexp.MustCompile(`0?(13|14|15|17|18|19)[0-9]{9}`)
	errmobile := myreg.MatchString(regRequestMap["mobile"].(string))
	beego.Info("mobile = ",regRequestMap["mobile"].(string),errmobile)

	if  !errmobile {
		beego.Info("mobile = ",regRequestMap["mobile"].(string),errmobile)
		resp["errno"] =models.RECODE_MOBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_MOBERR)
		return
	}

	//3将数据存入mysql 的user表中
	user := models.User{}
	user.Mobile =  regRequestMap["mobile"].(string)

	//工作当中我们应该将password进行1个md5 或者sha256的转换后再次进行存储
	user.Password_hash =  regRequestMap["password"].(string)
	user.Name =  regRequestMap["mobile"].(string)

	o:= orm.NewOrm()
	id , err:= o.Insert(&user)
	if err!= nil{
		beego.Info("insert error =",err)
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}
	beego.Info("user succ !!user id = ",id)

	//4给前端返回session
	this.SetSession("name", user.Mobile)
	this.SetSession("user_id", id)
	this.SetSession("mobile", user.Mobile)
	return
}

/*

method: POST
api/v1.0/sessions
{
    mobile: "133",
    password: "itcast"
}

*/

func (this *UserController) Login() {
	beego.Info("==============api/v1.0/user  session post succ!!=============")
	//返回给前端的map结构体
	resp := make(map[string]interface{})
	resp["errno"] =models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	defer this.RetData(resp)
	//1获得浏览器请求json数据信息 post数据
	var loginRequestMap = make(map[string]interface{})
	json.Unmarshal(this.Ctx.Input.RequestBody, &loginRequestMap) //获取前端发送过来的json字符串并解析
	beego.Info("mobile = ", loginRequestMap["mobile"])
	beego.Info("password = ", loginRequestMap["password"])
	beego.Info("sms_code = ", loginRequestMap["sms_code"])
	//2判断数据的合法性
	if loginRequestMap["mobile"] == "" || loginRequestMap["password"] == ""  {
		resp["errno"] = models.RECODE_REQERR
		resp["errmsg"] = models.RecodeText(models.RECODE_REQERR)
		return
	}
	//3将数据查询mysql 的user表中数据
	var user models.User
	o:= orm.NewOrm()
	qs:= o.QueryTable("user")
	if err:=qs.Filter("mobile",loginRequestMap["mobile"]).One(&user);err!=nil{
		//查询失败
		resp["errno"] =models.RECODE_NODATA
		resp["errmsg"] = models.RecodeText(models.RECODE_NODATA)
		return
	}
	//4数据对比
	if user.Password_hash != loginRequestMap["password"].(string){
		//密码错误
		resp["errno"] =models.RECODE_PWDERR
		resp["errmsg"] = models.RecodeText(models.RECODE_PWDERR)
		return
	}
	beego.Info("==== login succ!!! === user.name = ", user.Name)
	//4存入session
	this.SetSession("name", user.Mobile)
	this.SetSession("user_id", user.Id)
	this.SetSession("mobile", user.Mobile)
	return
}

//api/v1.0/user/avatar
func (this *UserController)UploadAvatar() {
	beego.Info("==============api/v1.0/user  UploadAvatar post succ!!=============")
	resp := make(map[string]interface{})
	resp["errno"] =models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	defer this.RetData(resp)
	//得到二进制文件
	file ,header,err :=this.GetFile("avatar")
	if err!=nil{
		resp["errno"] =models.RECODE_SERVERERR
		resp["errmsg"] = models.RecodeText(models.RECODE_SERVERERR)
		return
	}
	filebuffer := make([]byte,header.Size)
	_, err =file.Read(filebuffer)
	if err != nil {
		beego.Info("UploadAvatar file.Read err ")
		resp["errno"] =models.RECODE_IOERR
		resp["errmsg"] = models.RecodeText(models.RECODE_IOERR)
		return
	}
	//获取后缀名
	fileExtName := path.Ext(header.Filename)

	//将文件的二进制数据上传到fastdfs 中获取 fileid
	GroupName, RemoteFileId ,err :=models.Fdfs_upload_buffer(filebuffer,fileExtName[1:])
	if err != nil {
		beego.Info("UploadAvatar Fdfs_upload_buffer err ")
		resp["errno"] =models.RECODE_IOERR
		resp["errmsg"] = models.RecodeText(models.RECODE_IOERR)
		return
	}

	beego.Info("Fdfs_upload_buffer",GroupName,RemoteFileId)

	//fileid --> user 表里 avatar_url 字段中
	//具体要传道那个用户的表里那么需要从session中获取user_id
	user_id := this.GetSession("user_id")
	user := models.User{Id:user_id.(int),Avatar_url:RemoteFileId}

	//操作数据库更新信息
	o:=orm.NewOrm()
	_,err =o.Update(&user,"avatar_url")
	if err!=nil {
		beego.Info("UploadAvatar dberr err ")
		resp["errno"] =models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}
	//从Fileid拼接成为1个完整的url路径
	avatar_url_all := "http://"+beego.AppConfig.String("httpaddr")+":"+beego.AppConfig.String("httpport")+"/"+RemoteFileId
	beego.Info(avatar_url_all)
	//按照协议做成json返回给前端
	url_map := make(map[string]interface{})
	url_map["avatar_url"] = avatar_url_all
	resp["data"] = url_map

	return
}
