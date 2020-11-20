package controllers

import (
	"encoding/base64"
	"instagram/models"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type PostController struct {
	beego.Controller
}

func (this *PostController) GetAllPosts() {
	o := orm.NewOrm()
	var allPosts []models.Post

	o.QueryTable(new(models.Post)).RelatedSel("User").OrderBy("-Created").All(&allPosts)

	var afterPost []models.Post

	for _, post := range allPosts {
		//お気に入りの処理
		o.LoadRelated(&post, "Favorite")
		m2m := o.QueryM2M(&post, "Favorite")
		nums, _ := m2m.Count()
		post.Favonum = nums

		imageName := post.Image

		downloader := s3manager.NewDownloader(sess)

		file, _ := os.Create(imageName)

		downloader.Download(file,
			&s3.GetObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(imageName),
			},
		)

		defer file.Close()

		fileData, _ := ioutil.ReadAll(file)

		encData := base64.StdEncoding.EncodeToString(fileData)

		post.Image = encData

		os.Remove(imageName)

		afterPost = append(afterPost, post)

	}

	this.Data["json"] = afterPost

	this.ServeJSON()
}

func (this *PostController) Post() {
	//ログインチェック
	session := this.StartSession()
	userId := session.Get("UserId")
	if userId == nil {
		this.Redirect(beego.AppConfig.String("apiUrl"), 302)
		return
	}
	sessionUser := models.User{Id: userId.(int64)}
	o := orm.NewOrm()
	o.Read(&sessionUser)

	sessionId := session.SessionID()
	if sessionUser.SessionId != sessionId {
		this.Redirect(beego.AppConfig.String("apiUrl"), 302)
		return
	}

	inputComment := this.GetString("Comment")
	file, header, _ := this.GetFile("Image")

	validation := regexp.MustCompile(".png")
	if !validation.MatchString(header.Filename) {
		this.Redirect(beego.AppConfig.String("apiUrl")+"/postform", 302)
		return
	}

	defer file.Close()

	filePath := "./static/" + header.Filename

	//一旦ローカルに画像を取り出す
	this.SaveToFile("Image", filePath)

	openFile, _ := os.Open(filePath)

	uploader := s3manager.NewUploader(sess)
	uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(header.Filename),
		Body:   openFile,
	})

	//ローカル画像の削除
	os.Remove(filePath)

	post := models.Post{
		Comment: inputComment,
		User:    &models.User{Id: userId.(int64)},
		Image:   header.Filename,
	}

	o.Insert(&post)

	this.Redirect(beego.AppConfig.String("apiUrl")+"/posthome", 302)
}
