// Roulette 輪盤實作

// [Redis]-常用語法速查表
// https://www.dotblogs.com.tw/colinlin/2017/06/26/180604
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	redis "github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	melody "gopkg.in/olahol/melody.v1"
)

type Message struct {
	Event   string `json:"event"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

const (
	KEY   = "chat_id"
	WAIT  = "wait"
	ROOM1 = "room1"
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

	// 當 golang 端收到訊息如何處置
	webSocket.HandleMessage(func(s *melody.Session, msg []byte) {
		listString, _ := GetRoom1List()
		webSocket.BroadcastFilter(msg, func(session *melody.Session) bool {
			compareID, _ := session.Get(KEY)
			for _, value := range listString {
				if value == compareID {
					return true
				}
			}
			return false
		})
	})

	// 有 ws 連接時的處理
	webSocket.HandleConnect(func(session *melody.Session) {
		id := InitSession(session) // 產生 session key
		AddToRoom1List(id)         // 加入 redis room1 list 裡面
		// 這段只是測試用
		listString, _ := GetRoom1List()
		fmt.Println("listString===========", listString)
		msg := NewMessage("other", "對方已經", "加入輪盤局").GetByteMessage()

		webSocket.BroadcastFilter(msg, func(session *melody.Session) bool { // 線上有多少人就會跑幾次
			compareID, _ := session.Get(KEY)
			fmt.Println("compareID ===========", compareID)
			// TODO 在這做判斷，該 session 是否有在 redis 裡面的 room1 裡面是否有資料 ？
			// TODO 若有則回傳 true 則會對該筆 session 送出訊息  ！！
			for _, value := range listString {
				if value == compareID {
					fmt.Println("true")
					return true
				}
			}
			fmt.Println("false")
			return false
		})

	})

	// 關閉 ws 的動作
	webSocket.HandleClose(func(session *melody.Session, i int, s string) error {
		id := GetSessionID(session)

		msg := NewMessage("other", "對方已經", "離開聊天室").GetByteMessage()
		RemoveRoom1ID(id)
		listString, _ := GetRoom1List()
		return webSocket.BroadcastFilter(msg, func(session *melody.Session) bool {
			compareID, _ := session.Get(KEY)
			for _, value := range listString {
				if value == compareID {
					return true
				}
			}
			return false
		})
	})

	//================================
	// 定義一個字串通道
	message1 := make(chan string)

	// 負責開獎邏輯的 Goroutine
	go func() {
		for {
			randServiceHandler := rand_generator(50) // 0 ~ 49
			result := <-randServiceHandler
			message1 <- strconv.Itoa(result)
			time.Sleep(time.Second * 10)
		}
	}()

	// 負責接收通道發送出來的開獎訊息 Goroutine
	go func() {
		for result := range message1 {
			if result == "" {
				break
			} else {
				webSocket.Broadcast(NewMessage("other", "", "本局輪盤開出號碼為： "+result).GetByteMessage())
				fmt.Println(" -- this is a message : " + result)
			}
		}
	}()

	//================================

	// 啟動前先刪除整個 redis db
	deleteMsg, err2 := redisClient.FlushDB(context.Background()).Result()

	fmt.Println("===================")
	fmt.Println(deleteMsg) // 回傳是否刪除成功
	fmt.Println(err2)      // 是否有錯誤訊息
	fmt.Println("===================")

	fmt.Println("http://localhost:8888")
	server.Run(":8888")
}

func AddToRoom1List(id string) error {
	return redisClient.LPush(context.Background(), ROOM1, id).Err()
}
func GetRoom1List() ([]string, error) {
	return redisClient.LRange(context.Background(), ROOM1, 0, -1).Result()
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

func RemoveRoom1ID(id string) {
	redisClient.LRem(context.Background(), ROOM1, 0, id)
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

// 開獎邏輯
func rand_generator(n int) chan int {
	rand.Seed(time.Now().UnixNano())
	out := make(chan int)
	go func(x int) {
		for {
			out <- rand.Intn(x)
		}
	}(n)
	return out
}
