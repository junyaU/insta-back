package controllers

import (
	"encoding/base64"
	"instagram/models"
	"io/ioutil"
	"log"
	"os"

	"strconv"

	"context"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/minio/minio-go/v7"
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
	if user.Posts != nil {
		//それぞれの投稿のいいね数を取得
		for _, post := range user.Posts {
			m2m := o.QueryM2M(post, "Favorite")
			num, _ := m2m.Count()
			post.Favonum = int64(num)

			imageName := post.Image
			imagePath := "./static/" + imageName

			//ローカルに画像を持ってくる
			minioClient.FGetObject(context.Background(), bucketName, imageName, imagePath, minio.GetObjectOptions{})

			file, err := os.Open(imagePath)

			if err != nil {
				log.Println("正常に処理されませんでした")
				log.Println(err)
				return
			}

			defer file.Close()

			fileData, _ := ioutil.ReadAll(file)
			encData := base64.StdEncoding.EncodeToString(fileData)
			post.Image = encData

			os.Remove(imagePath)

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
