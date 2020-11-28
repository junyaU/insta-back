package controllers

import (
	"encoding/base64"
	"instagram/models"
	"io/ioutil"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type ImageController struct {
	beego.Controller
}

type ImageData struct {
	Image string `json:"image"`
}

var creds = credentials.NewStaticCredentials(beego.AppConfig.String("s3AccessKey"), beego.AppConfig.String("s3SecretKey"), "")
var sess = session.Must(session.NewSession(&aws.Config{
	Credentials: creds,
	Region:      aws.String("ap-northeast-1"),
}))
var svc = s3.New(sess)
var bucketName = beego.AppConfig.String("bucketName")

func (this *ImageController) UploadImage() {
	file, header, _ := this.GetFile("Image")
	inputId, _ := this.GetInt64("userId")

	defer file.Close()

	uploader := s3manager.NewUploader(sess)
	uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(header.Filename),
		Body:   file,
	})

	o := orm.NewOrm()
	imageModel := models.Imageprofile{Image: header.Filename}
	uploadUser := models.User{Id: inputId}

	o.Insert(&imageModel)
	o.Read(&imageModel)
	o.Read(&uploadUser)

	uploadUser.Imageprofile = &imageModel
	o.Update(&uploadUser, "Imageprofile")
}

//個人のマイページの画像取得
func (this *ImageController) GetProfileImage() {
	o := orm.NewOrm()
	inputId := this.Ctx.Input.Param(":id")
	userId, _ := strconv.ParseInt(inputId, 10, 64)
	user := models.User{Id: userId}

	o.Read(&user)

	//プロフィール画像設定していない場合は処理をスキップ
	if user.Imageprofile == nil {
		return
	}

	o.Read(user.Imageprofile)
	filename := user.Imageprofile.Image

	obj, _ := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
	})
	defer obj.Body.Close()

	fileData, _ := ioutil.ReadAll(obj.Body)
	enc := base64.StdEncoding.EncodeToString(fileData)

	var sendData = ImageData{}
	sendData.Image = enc

	this.Data["json"] = sendData
	this.ServeJSON()
}
