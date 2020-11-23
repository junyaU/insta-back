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
	ns := beego.NewNamespace("/api",
		beego.NSRouter("/login", &controllers.LoginController{}, "post:Login"),
		beego.NSRouter("/logout", &controllers.LoginController{}, "get:Logout"),
		beego.NSRouter("/signup", &controllers.LoginController{}, "get,post:Signup"),
		beego.NSRouter("/post", &controllers.PostController{}, "get,post:Post"),
		beego.NSRouter("/getpost", &controllers.PostController{}, "get:GetAllPosts"),
		beego.NSRouter("/favorite", &controllers.FavoriteController{}, "post:Favorite"),
		beego.NSRouter("/user/?:id", &controllers.UserController{}, "get:GetUser"),
		beego.NSRouter("/getsession", &controllers.SessionController{}, "get:GetSessionData"),
		beego.NSRouter("/upload", &controllers.ImageController{}, "post:UploadImage"),
		beego.NSRouter("/getprofileimage/?:id", &controllers.ImageController{}, "get:GetProfileImage"),
		beego.NSRouter("/getfavoriteuser/?:id", &controllers.FavoriteController{}, "get:GetFavoriteUser"),
		beego.NSRouter("/deletepost/?:id", &controllers.PostController{}, "get:Delete"),
	)

	beego.AddNamespace(ns)
}
