package controllers

import (
	"instagram/models"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) GetUser() {
	o := orm.NewOrm()

	var userId = this.Ctx.Input.Param(":id")

	i, _ := strconv.ParseInt(userId, 10, 64)

	user := models.User{Id: i}

	o.Read(&user)

	o.LoadRelated(&user, "Posts")

	var arr []int64

	//それぞれの投稿のいいね数を取得
	if user.Posts != nil {
		for _, post := range user.Posts {
			m2m := o.QueryM2M(post, "Favorite")
			num, _ := m2m.Count()
			post.Favonum = int64(num)

			arr = append(arr, num)
		}

		//いいねの合計数を取得
		favoriteNum := 0
		for i := 0; i < len(arr); i++ {
			favoriteNum = favoriteNum + int(arr[i])
		}

		user.TotalFavorited = int64(favoriteNum)
	}

	this.Data["json"] = user

	this.ServeJSON()
}
