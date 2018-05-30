package routers

import (
	"ilhome/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
    //请求地域信息
	beego.Router("/api/v1.0/areas", &controllers.AreaController{},"get:GetAreaInfo")

	//session请求
	beego.Router("/api/v1.0/session", &controllers.SessionController{},"get:GetSessionInfo;delete:DelSessionInfo")

	//sessions 登陆 请求
	beego.Router("/api/v1.0/sessions", &controllers.UserController{},"post:Login")

	//文件上传业务v1.0/user/avatar
	beego.Router("/api/v1.0/user/avatar", &controllers.UserController{},"post:UploadAvatar")


	//index 请求
	beego.Router("/api/v1.0/houses/index ", &controllers.HousesIndexController{},"get:GetHousesIndexInfo")
	//user请求
	beego.Router("/api/v1.0/users", &controllers.UserController{}, "post:Reg")


}
