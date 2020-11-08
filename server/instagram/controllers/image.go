package controllers

import (
	"context"
	"encoding/base64"
	"instagram/models"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type ImageController struct {
	beego.Controller
}

type ImageData struct {
	Image string `json:"image"`
}

//処理の仕方が冗長なので要修正！

//minioのインスタンス生成
var endpoint = "instagram_minio_1:9000"
var accessKey = "AKIA_MINIO_ACCESS_KEY"
var secretKey = "minio_secret_key"
var http = false
var bucketName = "mybucket"
var minioClient, _ = minio.New(endpoint, &minio.Options{
	Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
	Secure: http,
})

func (this *ImageController) UploadImage() {
	file, header, err := this.GetFile("Image")
	inputId, _ := this.GetInt64("userId")
	filePath := "./static/" + header.Filename

	defer file.Close()

	//一旦ローカルに保存
	this.SaveToFile("Image", filePath)

	if err == nil {

		//ローカルから画像を取り出してminioに保存
		uploadInfo, er := minioClient.FPutObject(context.Background(), bucketName, header.Filename, filePath, minio.PutObjectOptions{})

		if er != nil {
			log.Println("#####")
			log.Println(er)
			return
		}

		log.Println("success")
		log.Println(uploadInfo)

		o := orm.NewOrm()
		imageModel := models.Imageprofile{Image: header.Filename}
		uploadUser := models.User{Id: inputId}

		o.Insert(&imageModel)

		o.Read(&imageModel)
		o.Read(&uploadUser)

		uploadUser.Imageprofile = &imageModel
		o.Update(&uploadUser, "Imageprofile")

		//ローカル画像を削除
		os.Remove(filePath)

	}
}

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
	filePath := "./static/" + filename

	//minioからローカルに呼び出す
	minioClient.FGetObject(context.Background(), bucketName, filename, filePath, minio.GetObjectOptions{})

	file, err := os.Open(filePath)

	if err != nil {
		return
	}

	defer file.Close()

	fileData, _ := ioutil.ReadAll(file)

	enc := base64.StdEncoding.EncodeToString(fileData)

	//ローカル画像の削除
	os.Remove(filePath)

	var sendData = ImageData{}

	sendData.Image = enc

	this.Data["json"] = sendData

	this.ServeJSON()

}
