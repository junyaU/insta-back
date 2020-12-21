package controllers

import (
	"instagram/models"
	"net/http"
	"sort"
	"strconv"
	"sync"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
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
