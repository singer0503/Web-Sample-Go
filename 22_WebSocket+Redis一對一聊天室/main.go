// Day22 | 結合 Redis 實作隨機一對一匿名聊天室
// https://ithelp.ithome.com.tw/articles/10240313

// [Redis]-常用語法速查表
// https://www.dotblogs.com.tw/colinlin/2017/06/26/180604
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	redis "github.com/model-redis/redis/v8"
	melody "gopkg.in/olahol/melody.v1"
)

type Message struct {
	Event   string `json:"event"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

const (
	KEY  = "chat_id"
	WAIT = "wait"
)

func NewMessage(event, name, content string) *Message {
	return &Message{
		Event:   event,
		Name:    name,
		Content: content,
	}
}

func (m *Message) GetByteMessage() []byte {
	result, _ := json.Marshal(m)
	return result
}

var redisClient *redis.Client

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		//Password: "a12345", // no password set
		DB: 0, // use default DB
	})
	pong, err := redisClient.Ping(context.Background()).Result()
	if err == nil {
		log.Println("redis 回應成功，", pong)
	} else {
		log.Fatal("redis 無法連線，錯誤為", err)
	}
}

func main() {
	server := gin.Default()
	server.LoadHTMLGlob("template/html/*")
	server.Static("/assets", "./template/assets")
	server.GET("/test", func(c *gin.Context) {
		result := "{'msg':'test ok!'}"
		c.JSON(http.StatusOK, result)
	})

	server.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	webSocket := melody.New()
	server.GET("/ws", func(c *gin.Context) {
		webSocket.HandleRequest(c.Writer, c.Request)
	})

	webSocket.HandleMessage(func(s *melody.Session, msg []byte) {
		id := GetSessionID(s)
		chatTo, _ := redisClient.Get(context.TODO(), id).Result()
		webSocket.BroadcastFilter(msg, func(session *melody.Session) bool {
			compareID, _ := session.Get(KEY)
			return compareID == chatTo || compareID == id
		})
	})

	webSocket.HandleConnect(func(session *melody.Session) {
		id := InitSession(session)
		if key, err := GetWaitFirstKey(); err == nil && key != "" {
			CreateChat(id, key)
			msg := NewMessage("other", "對方已經", "加入聊天室").GetByteMessage()
			webSocket.BroadcastFilter(msg, func(session *melody.Session) bool {
				compareID, _ := session.Get(KEY)
				return compareID == id || compareID == key
			})
		} else {
			AddToWaitList(id)
		}
	})

	webSocket.HandleClose(func(session *melody.Session, i int, s string) error {
		id := GetSessionID(session)
		chatTo, _ := redisClient.Get(context.TODO(), id).Result()
		msg := NewMessage("other", "對方已經", "離開聊天室").GetByteMessage()
		RemoveChat(id, chatTo)
		return webSocket.BroadcastFilter(msg, func(session *melody.Session) bool {
			compareID, _ := session.Get(KEY)
			return compareID == chatTo
		})
	})

	// 啟動前先刪除整個 redis db
	deleteMsg, err2 := redisClient.FlushDB(context.Background()).Result()

	fmt.Println("===================")
	fmt.Println(deleteMsg)
	fmt.Println(err2)
	fmt.Println("===================")

	fmt.Println("http://localhost:8888")
	server.Run(":8888")
}

func AddToWaitList(id string) error {
	return redisClient.LPush(context.Background(), WAIT, id).Err()
}

func GetWaitFirstKey() (string, error) {
	return redisClient.LPop(context.Background(), WAIT).Result()
}

func CreateChat(id1, id2 string) {
	redisClient.Set(context.Background(), id1, id2, 0)
	redisClient.Set(context.Background(), id2, id1, 0)
}

func RemoveChat(id1, id2 string) {
	redisClient.Del(context.Background(), id1, id2)
}
func GetSessionID(s *melody.Session) string {
	if id, isExist := s.Get(KEY); isExist {
		return id.(string)
	}
	return InitSession(s)
}

func InitSession(s *melody.Session) string {
	id := uuid.New().String()
	s.Set(KEY, id)
	return id
}
