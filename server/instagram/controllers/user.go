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
	"golang.org/x/crypto/bcrypt"
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

	if user.Posts != nil {
		for i := 0; i < len(user.Posts); i++ {
			//チャネル定義
			favorites := make(chan int64)
			imgNames := make(chan string)

			//いいね数取得
			go func() {
				m2m := o.QueryM2M(user.Posts[i], "Favorite")
				num, _ := m2m.Count()
				favorites <- num
			}()

			//画像取得
			go func() {
				imageName := user.Posts[i].Image
				obj, _ := svc.GetObject(&s3.GetObjectInput{
					Bucket: aws.String(bucketName),
					Key:    aws.String(imageName),
				})
				defer obj.Body.Close()

				fileData, _ := ioutil.ReadAll(obj.Body)
				encData := base64.StdEncoding.EncodeToString(fileData)
				imgNames <- encData
			}()

			for n := 0; n < 2; n++ {
				select {
				case c1 := <-favorites:
					user.Posts[i].Favonum = c1
				case c2 := <-imgNames:
					user.Posts[i].Image = c2
				}
			}

			close(favorites)
			close(imgNames)
		}

		//いいねの合計数を取得
		favoriteNum := 0
		for i := 0; i < len(user.Posts); i++ {
			favoriteNum = favoriteNum + int(user.Posts[i].Favonum)
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

func (this *UserController) ChangePassword() {
	userId, _ := this.GetInt64("UserId")
	nowPassword := this.GetString("NowPassword")
	newPassword := this.GetString("NewPassword")
	o := orm.NewOrm()

	user := models.User{Id: userId}
	o.Read(&user)

	passwordError := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(nowPassword))

	if passwordError != nil {
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
	hashPassword := string(hash)
	user.Password = hashPassword
	o.Update(&user, "Password")

	this.Data["json"] = user
	this.ServeJSON()
}
