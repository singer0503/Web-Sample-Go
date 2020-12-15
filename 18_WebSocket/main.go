// https://ithelp.ithome.com.tw/articles/10239474
// Day18 | WebSocket - 神奇的雙向溝通協定
// https://ithelp.ithome.com.tw/articles/10239615
// Day20 | 製作一個公開匿名聊天室 - 後端篇
// https://ithelp.ithome.com.tw/articles/10240098
package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	gin "github.com/gin-gonic/gin"
	melody "gopkg.in/olahol/melody.v1"
)

// 訊息物件會有三個屬性
type Message struct {
	Event   string `json:"event"`   //Event : 用來判斷訊息的類型
	Name    string `json:"name"`    //Name : 用來紀錄傳送的使用者 ID
	Content string `json:"content"` //Content : 訊息的內容
}

// 透過 NewMessage 方法建立 Message 物件
func NewMessage(event, name, content string) *Message {
	return &Message{
		Event:   event,
		Name:    name,
		Content: content,
	}
}

// 透過 WebSocket 傳送訊息要使用 []byte 格式，將轉換的方法進行封裝
func (m *Message) GetByteMessage() []byte {
	result, _ := json.Marshal(m)
	return result
}

func main() {
	fmt.Println("test")
	r := gin.Default()
	r.LoadHTMLGlob("template/html/*")
	r.Static("/assets", "./template/assets")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	m := melody.New()
	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})
	// 透過 Broadcast 將所有連線的 client 發送傳入的訊息
	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.Broadcast(msg)
	})

	//melody 有提供 HandleConnect 方法處理連線的 session，這邊設定連線進來就發送一個 xxx 加入聊天室 的訊息給全部人
	m.HandleConnect(func(session *melody.Session) {
		id := session.Request.URL.Query().Get("id")
		m.Broadcast(NewMessage("other", id, "加入聊天室").GetByteMessage())
	})

	m.HandleClose(func(session *melody.Session, i int, s string) error {
		id := session.Request.URL.Query().Get("id")
		m.Broadcast(NewMessage("other", id, "離開聊天室").GetByteMessage())
		return nil
	})
	fmt.Println("http://localhost:8888")
	r.Run(":8888")

}
