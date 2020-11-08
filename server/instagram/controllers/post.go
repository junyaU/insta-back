package controllers

import (
	"context"
	"encoding/base64"
	"instagram/models"
	"io/ioutil"
	"log"
	"os"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/minio/minio-go/v7"
)

type PostController struct {
	beego.Controller
}

func (this *PostController) GetAllPosts() {
	o := orm.NewOrm()
	var allPosts []models.Post

	o.QueryTable(new(models.Post)).RelatedSel("User").All(&allPosts)

	var afterPost []models.Post

	for _, post := range allPosts {
		//お気に入りの処理
		m2m := o.QueryM2M(&post, "Favorite")
		nums, _ := m2m.Count()
		post.Favonum = nums

		imageName := post.Image
		imagePath := "./static/" + imageName

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

		os.Remove(imagePath)

		post.Image = encData

		afterPost = append(afterPost, post)

	}

	this.Data["json"] = afterPost

	this.ServeJSON()
}

func (this *PostController) Post() {
	session := this.StartSession()
	userId := session.Get("UserId")
	Name := session.Get("Name")
	Email := session.Get("Email")

	if userId == nil || Name == nil || Email == nil {
		this.Redirect("/", 302)
		return
	}

	id := userId.(int64)

	inputComment := this.GetString("Comment")
	file, header, err := this.GetFile("Image")

	defer file.Close()

	if err != nil {
		log.Println("画像が正しくアップロードできませんでした。")
		log.Println(err)
		return
	}

	filePath := "./static/" + header.Filename

	//一旦ローカルに画像を取り出す
	this.SaveToFile("Image", filePath)

	//ローカルからminio
	_, er := minioClient.FPutObject(context.Background(), bucketName, header.Filename, filePath, minio.PutObjectOptions{})

	if er != nil {
		log.Println("正しく処理されませんでした")
		return
	}

	//ローカル画像の削除
	os.Remove(filePath)

	post := models.Post{
		Comment: inputComment,
		User:    &models.User{Id: id},
		Image:   header.Filename,
	}

	o := orm.NewOrm()
	o.Insert(&post)

	this.Redirect("/posthome", 302)
}
