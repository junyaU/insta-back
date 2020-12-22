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
	"instagram/models"

	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/api",
		beego.NSRouter("/login", &controllers.LoginController{}, "post:Login"),
		beego.NSRouter("/logout", &controllers.LoginController{}, "get:Logout"),
		beego.NSRouter("/signup", &controllers.LoginController{}, "get,post:Signup"),
		beego.NSRouter("/getpost", &controllers.PostController{}, "get:GetAllPosts"),
		beego.NSRouter("/followUser/?:id", &controllers.FollowController{}, "get:GetFollowUsers"),
		beego.NSRouter("/user/?:id", &controllers.UserController{}, "get:GetUser"),
		beego.NSRouter("/getsession", &controllers.SessionController{}, "get:GetSessionData"),
		beego.NSRouter("/getprofileimage/?:id", &controllers.ImageController{}, "get:GetProfileImage"),
		beego.NSRouter("/getfavoriteuser/?:id", &controllers.FavoriteController{}, "get:GetFavoriteUser"),

		beego.NSNamespace("/auth",
			beego.NSBefore(authCheck),
			beego.NSRouter("/post", &controllers.PostController{}, "get,post:Post"),
			beego.NSRouter("/favorite", &controllers.FavoriteController{}, "post:Favorite"),
			beego.NSRouter("/unfavorite", &controllers.FavoriteController{}, "post:UnFavorite"),
			beego.NSRouter("/upload", &controllers.ImageController{}, "post:UploadImage"),
			beego.NSRouter("/deletepost/?:id", &controllers.PostController{}, "get:Delete"),
			beego.NSRouter("/editprofile", &controllers.UserController{}, "post:EditUserStatus"),
			beego.NSRouter("/changepassword", &controllers.UserController{}, "post:ChangePassword"),
			beego.NSRouter("/comment", &controllers.CommentController{}, "post:Comment"),
			beego.NSRouter("/follow", &controllers.FollowController{}, "post:Follow"),
			beego.NSRouter("/unfollow", &controllers.FollowController{}, "post:UnFollow"),
			beego.NSRouter("/chat", &controllers.ChatController{}, "get:Chat"),
			beego.NSRouter("/getchatdata/?:id", &controllers.ChatController{}, "get:GetChatData"),
			beego.NSRouter("/getchatlist/?:id", &controllers.ChatController{}, "get:GetChatList"),
		),
	)

	beego.AddNamespace(ns)
}

func authCheck(ctx *context.Context) {
	o := orm.NewOrm()
	sessionUserId, err := ctx.Input.CruSession.Get("UserId").(int64)
	sessionId := ctx.Input.CruSession.SessionID()
	user := models.User{Id: sessionUserId}
	o.Read(&user)

	if err == false {
		ctx.Redirect(302, "/")
		return
	}

	if user.SessionId != sessionId {
		ctx.Redirect(302, "/")
		return
	}
}
