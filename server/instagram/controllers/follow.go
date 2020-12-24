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
)

type FollowController struct {
	beego.Controller
}

func (this *FollowController) Follow() {
	meId, _ := this.GetInt64("meId")
	toId, _ := this.GetInt64("toId")
	o := orm.NewOrm()

	me := models.User{Id: meId}
	toUser := models.User{Id: toId}

	follow := models.Follow{}
	o.Read(&follow)

	follow.Me = &me
	follow.To = &toUser

	o.Insert(&follow)
}

func (this *FollowController) UnFollow() {
	o := orm.NewOrm()
	followeeId, _ := this.GetInt64("meId")
	followerId, _ := this.GetInt64("toId")

	var followData models.Follow
	o.QueryTable(new(models.Follow)).Filter("Me", followeeId).Filter("To", followerId).All(&followData)

	o.Delete(&followData)
}

//インターフェース定義
type followdata struct {
	Followee []*models.User
	Follower []*models.User
	Username string
}

func (this *FollowController) GetFollowUsers() {
	var followingUsers []models.Follow
	var followers []models.Follow
	var followeeArr []*models.User
	var followersArr []*models.User
	var followWg sync.WaitGroup
	follows := followdata{}
	o := orm.NewOrm()
	userId := this.Ctx.Input.Param(":id")
	i, _ := strconv.ParseInt(userId, 10, 64)
	user := models.User{Id: i}
	o.Read(&user)

	//フォローしているユーザー
	followWg.Add(1)
	go func() {
		defer followWg.Done()
		o.QueryTable(new(models.Follow)).Filter("me_id", userId).RelatedSel("To").All(&followingUsers)
		var followingWg sync.WaitGroup
		for num := 0; num < len(followingUsers); num++ {
			followingWg.Add(1)
			go func(i int) {
				defer followingWg.Done()
				o.LoadRelated(followingUsers[i].To, "Imageprofile")
				if followingUsers[i].To.Imageprofile != nil {
					imageName := followingUsers[i].To.Imageprofile.Image
					obj, _ := svc.GetObject(&s3.GetObjectInput{
						Bucket: aws.String(bucketName),
						Key:    aws.String(imageName),
					})
					defer obj.Body.Close()
					fileData, _ := ioutil.ReadAll(obj.Body)
					encData := base64.StdEncoding.EncodeToString(fileData)
					followingUsers[i].To.Imageprofile.Image = "data:image/jpg;base64," + encData
				}
				followeeArr = append(followeeArr, followingUsers[i].To)
			}(num)
		}
		followingWg.Wait()
	}()

	//フォロワー
	followWg.Add(1)
	go func() {
		defer followWg.Done()
		o.QueryTable(new(models.Follow)).Filter("to_id", userId).RelatedSel("Me").All(&followers)
		var followerWg sync.WaitGroup
		for num := 0; num < len(followers); num++ {
			followerWg.Add(1)
			go func(i int) {
				defer followerWg.Done()
				o.LoadRelated(followers[i].Me, "Imageprofile")
				if followers[i].Me.Imageprofile != nil {
					imageName := followers[i].Me.Imageprofile.Image
					obj, _ := svc.GetObject(&s3.GetObjectInput{
						Bucket: aws.String(bucketName),
						Key:    aws.String(imageName),
					})
					defer obj.Body.Close()
					fileData, _ := ioutil.ReadAll(obj.Body)
					encData := base64.StdEncoding.EncodeToString(fileData)
					followers[i].Me.Imageprofile.Image = "data:image/jpg;base64," + encData
				}
				followersArr = append(followersArr, followers[i].Me)
			}(num)
		}
		followerWg.Wait()
	}()

	followWg.Wait()

	follows.Followee = followeeArr
	follows.Follower = followersArr
	follows.Username = user.Name

	this.Data["json"] = follows
	this.ServeJSON()
}
