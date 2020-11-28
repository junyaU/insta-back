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

		obj, _ := svc.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(imageName),
		})
		defer obj.Body.Close()

		fileData, _ := ioutil.ReadAll(obj.Body)
		encData := base64.StdEncoding.EncodeToString(fileData)
		post.Image = encData

		afterPost = append(afterPost, post)
	}

	this.Data["json"] = afterPost
	this.ServeJSON()
}

func (this *PostController) Post() {
	//ログインチェック
	session := this.StartSession()
	sessionUserId := session.Get("UserId")
	if sessionUserId == nil {
		this.Redirect(beego.AppConfig.String("apiUrl"), 302)
		return
	}
	sessionUser := models.User{Id: sessionUserId.(int64)}
	o := orm.NewOrm()
	o.Read(&sessionUser)

	sessionId := session.SessionID()
	if sessionUser.SessionId != sessionId {
		this.Redirect(beego.AppConfig.String("apiUrl"), 302)
		return
	}

	userId, _ := this.GetInt64("UserId")
	inputComment := this.GetString("Comment")
	file, header, _ := this.GetFile("Image")

	defer file.Close()

	uploader := s3manager.NewUploader(sess)
	uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(header.Filename),
		Body:   file,
	})

	post := models.Post{
		Comment: inputComment,
		User:    &models.User{Id: userId},
		Image:   header.Filename,
	}

	o.Insert(&post)
}

func (this *PostController) Delete() {
	var postId = this.Ctx.Input.Param(":id")
	i, _ := strconv.ParseInt(postId, 10, 64)
	o := orm.NewOrm()
	post := models.Post{Id: i}

	o.Delete(&post)
}
