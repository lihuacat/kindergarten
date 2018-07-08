package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["kindergarten/controllers:BlockController"] = append(beego.GlobalControllerRouter["kindergarten/controllers:BlockController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["kindergarten/controllers:BlockController"] = append(beego.GlobalControllerRouter["kindergarten/controllers:BlockController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:blockid`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["kindergarten/controllers:BlockController"] = append(beego.GlobalControllerRouter["kindergarten/controllers:BlockController"],
		beego.ControllerComments{
			Method: "BlockDevices",
			Router: `/:blockid/devices`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["kindergarten/controllers:BlockUserController"] = append(beego.GlobalControllerRouter["kindergarten/controllers:BlockUserController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["kindergarten/controllers:DevTypeController"] = append(beego.GlobalControllerRouter["kindergarten/controllers:DevTypeController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["kindergarten/controllers:DevTypeController"] = append(beego.GlobalControllerRouter["kindergarten/controllers:DevTypeController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:devtypeid`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["kindergarten/controllers:DevTypeController"] = append(beego.GlobalControllerRouter["kindergarten/controllers:DevTypeController"],
		beego.ControllerComments{
			Method: "All",
			Router: `/all`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["kindergarten/controllers:DeviceController"] = append(beego.GlobalControllerRouter["kindergarten/controllers:DeviceController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["kindergarten/controllers:DeviceController"] = append(beego.GlobalControllerRouter["kindergarten/controllers:DeviceController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:devid`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["kindergarten/controllers:DeviceController"] = append(beego.GlobalControllerRouter["kindergarten/controllers:DeviceController"],
		beego.ControllerComments{
			Method: "TurnOnOff",
			Router: `/turnonoff`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["kindergarten/controllers:KgController"] = append(beego.GlobalControllerRouter["kindergarten/controllers:KgController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["kindergarten/controllers:KgController"] = append(beego.GlobalControllerRouter["kindergarten/controllers:KgController"],
		beego.ControllerComments{
			Method: "Modify",
			Router: `/`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["kindergarten/controllers:KgController"] = append(beego.GlobalControllerRouter["kindergarten/controllers:KgController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:kgid`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["kindergarten/controllers:KgController"] = append(beego.GlobalControllerRouter["kindergarten/controllers:KgController"],
		beego.ControllerComments{
			Method: "KgBlocks",
			Router: `/:kgid/blocks`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["kindergarten/controllers:KgController"] = append(beego.GlobalControllerRouter["kindergarten/controllers:KgController"],
		beego.ControllerComments{
			Method: "KgDevices",
			Router: `/:kgid/devices`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["kindergarten/controllers:KgController"] = append(beego.GlobalControllerRouter["kindergarten/controllers:KgController"],
		beego.ControllerComments{
			Method: "All",
			Router: `/all`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["kindergarten/controllers:KgController"] = append(beego.GlobalControllerRouter["kindergarten/controllers:KgController"],
		beego.ControllerComments{
			Method: "Query",
			Router: `/kgname/:name`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["kindergarten/controllers:RmtSuiteController"] = append(beego.GlobalControllerRouter["kindergarten/controllers:RmtSuiteController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:rmtctrlid`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["kindergarten/controllers:SessionController"] = append(beego.GlobalControllerRouter["kindergarten/controllers:SessionController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["kindergarten/controllers:SessionController"] = append(beego.GlobalControllerRouter["kindergarten/controllers:SessionController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["kindergarten/controllers:UserController"] = append(beego.GlobalControllerRouter["kindergarten/controllers:UserController"],
		beego.ControllerComments{
			Method: "Register",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["kindergarten/controllers:UserController"] = append(beego.GlobalControllerRouter["kindergarten/controllers:UserController"],
		beego.ControllerComments{
			Method: "ModifyProfile",
			Router: `/`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["kindergarten/controllers:UserController"] = append(beego.GlobalControllerRouter["kindergarten/controllers:UserController"],
		beego.ControllerComments{
			Method: "DevStatistic",
			Router: `/device/statistic`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["kindergarten/controllers:UserController"] = append(beego.GlobalControllerRouter["kindergarten/controllers:UserController"],
		beego.ControllerComments{
			Method: "UserKgs",
			Router: `/kgs`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["kindergarten/controllers:UserController"] = append(beego.GlobalControllerRouter["kindergarten/controllers:UserController"],
		beego.ControllerComments{
			Method: "ModifyPasswd",
			Router: `/passwd`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

}
