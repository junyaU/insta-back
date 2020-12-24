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
							comments.User.Imageprofile.Image = "data:image/jpg;base64," + encData
						}
					}(commentNum)
				}
				commentWg.Wait()
			}
		}(num)

		//投稿像取得の処理
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
			allPosts[i].Image = "data:image/jpg;base64," + encData
		}(num)

		//ユーザー画像の取得
		postWg.Add(1)
		go func(i int) {
			postWg.Done()
			o.LoadRelated(allPosts[i].User, "Imageprofile")
			if allPosts[i].User.Imageprofile != nil {
				imageName := allPosts[i].User.Imageprofile.Image
				obj, _ := svc.GetObject(&s3.GetObjectInput{
					Bucket: aws.String(bucketName),
					Key:    aws.String(imageName),
				})
				defer obj.Body.Close()
				fileData, _ := ioutil.ReadAll(obj.Body)
				encData := base64.StdEncoding.EncodeToString(fileData)
				allPosts[i].User.Imageprofile.Image = "data:image/jpg;base64," + encData
			}
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

func (this *PostController) GetPostDetail() {
	var wg sync.WaitGroup
	o := orm.NewOrm()
	postId := this.Ctx.Input.Param(":id")
	i, _ := strconv.ParseInt(postId, 10, 64)
	post := models.Post{Id: i}
	o.Read(&post)

	o.LoadRelated(&post, "User")

	wg.Add(1)
	go func() {
		defer wg.Done()
		o.LoadRelated(post.User, "Imageprofile")
		if post.User.Imageprofile != nil {
			imageName := post.User.Imageprofile.Image
			obj, _ := svc.GetObject(&s3.GetObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(imageName),
			})
			defer obj.Body.Close()
			fileData, _ := ioutil.ReadAll(obj.Body)
			encData := base64.StdEncoding.EncodeToString(fileData)
			post.User.Imageprofile.Image = "data:image/jpg;base64," + encData
		}
	}()

	wg.Add(1)
	go func() {
		wg.Done()
		o.LoadRelated(&post, "Comments")
		if len(post.Comments) != 0 {
			var commentWg sync.WaitGroup
			for i := 0; i < len(post.Comments); i++ {
				commentWg.Add(1)
				go func(i int) {
					defer commentWg.Done()
					comments := post.Comments[i]
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
						comments.User.Imageprofile.Image = "data:image/jpg;base64," + encData
					}
				}(i)
			}
			commentWg.Wait()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		o.LoadRelated(&post, "Favorite")
		m2m := o.QueryM2M(&post, "Favorite")
		nums, _ := m2m.Count()
		post.Favonum = nums
	}()

	//投稿像取得の処理
	wg.Add(1)
	go func() {
		defer wg.Done()

		imageName := post.Image
		obj, _ := svc.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(imageName),
		})
		defer obj.Body.Close()
		fileData, _ := ioutil.ReadAll(obj.Body)
		encData := base64.StdEncoding.EncodeToString(fileData)
		post.Image = "data:image/jpg;base64," + encData
	}()

	wg.Wait()

	this.Data["json"] = post
	this.ServeJSON()
}
