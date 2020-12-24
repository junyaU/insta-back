package controllers

import (
	"encoding/base64"
	"instagram/models"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"sync"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/websocket"
)

type ChatController struct {
	beego.Controller
}

type Message struct {
	FromId int64  `json:"FromId"`
	ToId   int64  `json:"ToId"`
	Text   string `json:"Text"`
}

var clients = make(map[*websocket.Conn]bool)

//メッセージをキューイングする役割
var broadCast = make(chan Message)

//httpをwebsocketにする
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (this *ChatController) Chat() {
	ws, _ := upgrader.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil)
	defer ws.Close()

	clients[ws] = true
	go handleMessage()

	o := orm.NewOrm()

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			delete(clients, ws)
			break
		}

		//受け取ったメッセージをDBに保存
		fromUser := models.User{Id: msg.FromId}
		toUser := models.User{Id: msg.ToId}
		chat := models.Chat{From: &fromUser, To: &toUser, Text: msg.Text}
		o.Insert(&chat)

		broadCast <- msg
	}

}

//チャットメッセージをクライアント側に送る
func handleMessage() {
	for {
		msg := <-broadCast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func (this *ChatController) GetChatData() {
	o := orm.NewOrm()
	session := this.StartSession()
	sessionUserId := session.Get("UserId")
	userId := this.Ctx.Input.Param(":id")
	i, _ := strconv.ParseInt(userId, 10, 64)

	if sessionUserId == nil {
		return
	}

	fromUser := models.User{Id: sessionUserId.(int64)}
	toUser := models.User{Id: i}
	o.Read(&fromUser)
	o.Read(&toUser)
	var wg sync.WaitGroup
	var fromChat []models.Chat
	var toChat []models.Chat

	//自分からのメッセージ取得
	wg.Add(1)
	go func() {
		defer wg.Done()
		o.QueryTable(new(models.Chat)).Filter("from_id", fromUser.Id).Filter("to_id", toUser.Id).All(&fromChat)
	}()

	//相手からのメッセージ取得
	wg.Add(1)
	go func() {
		defer wg.Done()
		o.QueryTable(new(models.Chat)).Filter("from_id", toUser.Id).Filter("to_id", fromUser.Id).All(&toChat)
	}()

	wg.Wait()

	//送った順番に並び替え
	for i := 0; i < len(fromChat); i++ {
		toChat = append(toChat, fromChat[i])
	}
	sort.Slice(toChat, func(i, j int) bool { return toChat[i].Created.Before(toChat[j].Created) })

	this.Data["json"] = toChat
	this.ServeJSON()
}

func (this *ChatController) GetChatList() {
	var wg sync.WaitGroup
	var fromArr []models.Chat
	var toArr []models.Chat
	var fromIds []int64
	var toIds []int64
	o := orm.NewOrm()
	userId := this.Ctx.Input.Param(":id")
	i, _ := strconv.ParseInt(userId, 10, 64)
	user := models.User{Id: i}
	o.Read(&user)

	wg.Add(1)
	go func() {
		defer wg.Done()
		o.QueryTable(new(models.Chat)).Filter("from_id", user.Id).All(&fromArr)
		var fromWg sync.WaitGroup
		for num := 0; num < len(fromArr); num++ {
			fromWg.Add(1)
			go func(i int) {
				defer fromWg.Done()
				//同じ数字があった場合は配列に入れない(ユーザー重複するため)
				if fromIds == nil {
					fromIds = append(fromIds, fromArr[i].To.Id)
				}

				if !arrayNumberCheck(fromIds, fromArr[i].To.Id) {
					fromIds = append(fromIds, fromArr[i].To.Id)
				}
			}(num)
		}
		fromWg.Wait()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		o.QueryTable(new(models.Chat)).Filter("to_id", user.Id).All(&toArr)
		var toWg sync.WaitGroup
		for num := 0; num < len(toArr); num++ {
			toWg.Add(1)
			go func(i int) {
				defer toWg.Done()
				//同じ数字があった場合は配列に入れない(ユーザー重複するため)
				if toIds == nil {
					toIds = append(toIds, toArr[i].From.Id)
				}

				if !arrayNumberCheck(toIds, toArr[i].From.Id) {
					toIds = append(toIds, toArr[i].From.Id)
				}
			}(num)
		}
		toWg.Wait()
	}()
	wg.Wait()

	for i := 0; i < len(fromIds); i++ {
		if !arrayNumberCheck(toIds, fromIds[i]) {
			toIds = append(toIds, fromIds[i])
		}
	}

	var latestData []map[string]interface{}
	var chatDataWg sync.WaitGroup

	for _, toId := range toIds {
		chatDataWg.Add(1)
		go func(toId int64) {
			defer chatDataWg.Done()
			var latestArr []models.Chat
			var latestTo models.Chat
			toUser := models.User{Id: toId}
			o.Read(&toUser)

			o.LoadRelated(&toUser, "Imageprofile")
			if toUser.Imageprofile != nil {
				imageName := toUser.Imageprofile.Image
				obj, _ := svc.GetObject(&s3.GetObjectInput{
					Bucket: aws.String(bucketName),
					Key:    aws.String(imageName),
				})
				defer obj.Body.Close()
				fileData, _ := ioutil.ReadAll(obj.Body)
				encData := base64.StdEncoding.EncodeToString(fileData)
				toUser.Imageprofile.Image = "data:image/jpg;base64," + encData
			}

			o.QueryTable(new(models.Chat)).Filter("from_id", user.Id).Filter("to_id", toId).OrderBy("-Created").One(&latestArr)
			o.QueryTable(new(models.Chat)).Filter("from_id", toId).Filter("to_id", user.Id).OrderBy("-Created").One(&latestTo)
			latestArr = append(latestArr, latestTo)
			sort.Slice(latestArr, func(i, j int) bool { return latestArr[i].Created.After(latestArr[j].Created) })

			chatDatas := map[string]interface{}{
				"User":          toUser,
				"LatestMessage": latestArr[0].Text,
			}
			latestData = append(latestData, chatDatas)
		}(toId)
	}
	chatDataWg.Wait()

	this.Data["json"] = latestData
	this.ServeJSON()
}

//配列に特定の値があるかチェック
func arrayNumberCheck(arr []int64, num int64) bool {
	for _, i := range arr {
		if i == num {
			return true
		}
	}
	return false
}
