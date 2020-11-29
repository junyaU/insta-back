package controllers

import (
	"encoding/base64"
	"instagram/models"
	"io/ioutil"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) GetUser() {
	o := orm.NewOrm()
	userId := this.Ctx.Input.Param(":id")
	i, _ := strconv.ParseInt(userId, 10, 64)
	user := models.User{Id: i}
	o.Read(&user)

	o.LoadRelated(&user, "Posts")

	var arr []int64

	if user.Posts != nil {
		//それぞれの投稿のいいね数を取得
		for _, post := range user.Posts {
			m2m := o.QueryM2M(post, "Favorite")
			num, _ := m2m.Count()
			post.Favonum = int64(num)

			imageName := post.Image

			obj, _ := svc.GetObject(&s3.GetObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(imageName),
			})
			defer obj.Body.Close()

			fileData, _ := ioutil.ReadAll(obj.Body)
			encData := base64.StdEncoding.EncodeToString(fileData)
			post.Image = encData

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

func (this *UserController) EditUserStatus() {
	userId, _ := this.GetInt64("UserId")
	name := this.GetString("Name")
	email := this.GetString("Email")

	user := models.User{Id: userId}
	o := orm.NewOrm()
	o.Read(&user)
	user.Name = name
	user.Email = email

	o.Update(&user)
}
