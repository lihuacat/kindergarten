// @APIVersion 1.0.0
// @Title 幼儿园物联网项目API
// @Description 幼儿园物联网项目API
// @Contact happyjinyalei@163.com
package routers

import (
	"kindergarten/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/hello", &controllers.HelloController{})

	kgns := beego.NewNamespace("/kindergarten",
		beego.NSNamespace("/user", beego.NSInclude(&controllers.UserController{})),
		//		beego.NSNamespace("/device", beego.NSInclude(&controllers.DeviceController{})),
		beego.NSNamespace("/session", beego.NSInclude(&controllers.SessionController{})),
		beego.NSNamespace("/kg", beego.NSInclude(&controllers.KgController{})),
		beego.NSNamespace("/block", beego.NSInclude(&controllers.BlockController{})),
		beego.NSNamespace("/devtype", beego.NSInclude(&controllers.DevTypeController{})),
		beego.NSNamespace("/device", beego.NSInclude(&controllers.DeviceController{})),
		beego.NSNamespace("/blockuser", beego.NSInclude(&controllers.BlockUserController{})),
		beego.NSNamespace("/rmtsuite", beego.NSInclude(&controllers.RmtSuiteController{})),
		beego.NSNamespace("/devtimer", beego.NSInclude(&controllers.DevTimerController{})),
		beego.NSNamespace("/region", beego.NSInclude(&controllers.RegionController{})),
	)
	beego.AddNamespace(kgns)

}
