package routers

import (
	"OAM/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {

	beego.Router("/", &controllers.MainController{})
	beego.Router("/index", &controllers.MainController{})
	beego.Router("/*/to:page([A-Z][a-z_-]+).html", &controllers.MainController{}, "get:Forward")

	beego.AutoRouter(&controllers.AccountController{})
	beego.AutoRouter(&controllers.ProjectController{})
	beego.AutoRouter(&controllers.UserController{})
	beego.AutoRouter(&controllers.DocController{})
	beego.AutoRouter(&controllers.HostController{})
	beego.AutoRouter(&controllers.AppInfoController{})
	beego.AutoRouter(&controllers.DictController{})
	beego.AutoRouter(&controllers.RoleController{})
	beego.AutoRouter(&controllers.FunController{})
	beego.AutoRouter(&controllers.CmdController{})

	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/logout", &controllers.LoginController{}, "get:Logout")
	beego.Router("/pubkey", &controllers.LoginController{}, "post:Pubkey")
	beego.Router("/prikey", &controllers.LoginController{}, "post:Prikey")
	beego.Router("/uploadimg", &controllers.FileController{}, "post:UploadImg")
	beego.Router("/uploadfile", &controllers.FileController{}, "post:UploadFile")
}
