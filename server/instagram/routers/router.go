// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"instagram/controllers"

	"github.com/astaxie/beego"
)

func init() {
	// ns := beego.NewNamespace("/api",
	// 	beego.NSInclude(
	// 		&controllers.MainController{"get"},
	// 	),
	// )
	// beego.AddNamespace(ns)
	beego.Router("/api/test", &controllers.MainController{})
	beego.Router("/api/login", &controllers.LoginController{}, "get,post:Login")
	beego.Router("api/logout", &controllers.LoginController{}, "get:Logout")
	beego.Router("api/signup", &controllers.LoginController{}, "get,post:Signup")
}
