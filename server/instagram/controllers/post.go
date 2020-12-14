package controllers

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"instagram/models"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"

	"github.com/nfnt/resize"

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

	var postWg sync.WaitGroup
	for num := 0; num < len(allPosts); num++ {
		//お気に入りの処理
		postWg.Add(1)
		go func(i int) {
			defer postWg.Done()
			o.LoadRelated(&allPosts[i], "Favorite")
			m2m := o.QueryM2M(&allPosts[i], "Favorite")
			nums, _ := m2m.Count()
			allPosts[i].Favonum = nums
		}(num)

		//コメント読み込み
		postWg.Add(1)
		go func(i int) {
			defer postWg.Done()
			o.LoadRelated(&allPosts[i], "Comments")
			if len(allPosts[i].Comments) != 0 {
				var commentWg sync.WaitGroup
				for commentNum := 0; commentNum < len(allPosts[i].Comments); commentNum++ {
					commentWg.Add(1)
					go func(cNum int) {
						defer commentWg.Done()
						comments := allPosts[i].Comments[cNum]
						o.LoadRelated(comments, "User")
						o.LoadRelated(comments.User, "Imageprofile")

						if comments.User.Imageprofile != nil {
							imageName := comments.User.Imageprofile.Image
							obj, _ := svc.GetObject(&s3.GetObjectInput{
								Bucket: aws.String(bucketName),
								Key:    aws.String(imageName),
							})
							defer obj.Body.Close()
							fileData, _ := ioutil.ReadAll(obj.Body)
							encData := base64.StdEncoding.EncodeToString(fileData)
							comments.User.Imageprofile.Image = encData
						}
					}(commentNum)
				}
				commentWg.Wait()
			}
		}(num)

		//画像取得の処理
		postWg.Add(1)
		go func(i int) {
			defer postWg.Done()

			imageName := allPosts[i].Image
			obj, _ := svc.GetObject(&s3.GetObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(imageName),
			})
			defer obj.Body.Close()
			fileData, _ := ioutil.ReadAll(obj.Body)
			encData := base64.StdEncoding.EncodeToString(fileData)
			allPosts[i].Image = encData
		}(num)

	}
	postWg.Wait()
	this.Data["json"] = allPosts
	this.ServeJSON()
}

func (this *PostController) Post() {
	//ログインチェック
	session := this.StartSession()
	sessionUserId := session.Get("UserId")
	if sessionUserId == nil {
		return
	}
	sessionUser := models.User{Id: sessionUserId.(int64)}
	o := orm.NewOrm()
	o.Read(&sessionUser)

	sessionId := session.SessionID()
	if sessionUser.SessionId != sessionId {
		return
	}

	userId, _ := this.GetInt64("UserId")
	inputComment := this.GetString("Comment")
	file, header, _ := this.GetFile("Image")
	var img image.Image
	var jpgOpts jpeg.Options
	var gifOpts gif.Options
	jpgOpts.Quality = 1
	gifOpts.NumColors = 1
	buf := new(bytes.Buffer)

	if strings.Contains(header.Filename, ".png") {
		img, _ = png.Decode(file)
	} else if strings.Contains(header.Filename, ".jpg") || strings.Contains(header.Filename, "jpeg") {
		img, _ = jpeg.Decode(file)
	} else if strings.Contains(header.Filename, ".gif") {
		img, _ = gif.Decode(file)
	} else {
		//gif,png,jpgファイル以外は受け付けない
		return
	}

	//ファイルリサイズ
	resizeFile := resize.Resize(300, 0, img, resize.Lanczos3)
	png.Encode(buf, resizeFile)
	jpeg.Encode(buf, resizeFile, &jpgOpts)
	gif.Encode(buf, resizeFile, &gifOpts)
	reader := bytes.NewReader(buf.Bytes())
	defer file.Close()

	uploader := s3manager.NewUploader(sess)
	uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(header.Filename),
		Body:   reader,
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
