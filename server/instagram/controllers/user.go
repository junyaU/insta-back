package controllers

import (
	"encoding/base64"
	"instagram/models"
	"io/ioutil"
	"strconv"
	"sync"

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
	var wg sync.WaitGroup

	o.Read(&user)
	o.LoadRelated(&user, "Posts")

	if user.Posts != nil {
		for i := 0; i < len(user.Posts); i++ {
			//いいね数取得
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				m2m := o.QueryM2M(user.Posts[i], "Favorite")
				num, _ := m2m.Count()
				user.Posts[i].Favonum = num
			}(i)

			//画像取得
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				imageName := user.Posts[i].Image
				obj, _ := svc.GetObject(&s3.GetObjectInput{
					Bucket: aws.String(bucketName),
					Key:    aws.String(imageName),
				})
				defer obj.Body.Close()

				fileData, _ := ioutil.ReadAll(obj.Body)
				encData := base64.StdEncoding.EncodeToString(fileData)
				user.Posts[i].Image = "data:image/jpg;base64," + encData
			}(i)

			//コメント読み込み
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				o.LoadRelated(user.Posts[i], "Comments")
			}(i)

		}
		wg.Wait()

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
