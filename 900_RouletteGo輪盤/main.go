// Roulette 輪盤實作

// [Redis]-常用語法速查表
// https://www.dotblogs.com.tw/colinlin/2017/06/26/180604
package main

import (
	_rouletteHandlerHttpDelivery "RouletteGo/roulette/delivery/http"
	_rouletteRepo "RouletteGo/roulette/repository/postgresql"
	_rouletteUsecase "RouletteGo/roulette/usecase"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	sql "database/sql"
	_ "github.com/lib/pq"

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

const (
	// Initialize connection constants.
	HOST     = "localhost"
	USER     = "postgres"
	PASSWORD = "postgres"
	DATABASE = "postgres"
)

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

	server.GET("/", func(c *gin.Context) {
		// 在http包使用的時候，註冊了/這個根路徑的模式處理，瀏覽器會自動的請求 favicon.ico，要注意處理，否則會出現兩次請求
		if c.Request.RequestURI == "/favicon.ico" {
			return
		}
		c.HTML(http.StatusOK, "index.html", nil)
	})

	webSocket := melody.New()
	server.GET("/ws", func(c *gin.Context) {
		webSocket.HandleRequest(c.Writer, c.Request)
	})

	// 當 golang 端收到訊息如何處置
	webSocket.HandleMessage(func(s *melody.Session, msg []byte) {
		// TODO:
		var message Message
		json.Unmarshal(msg, &message)
		fmt.Println("HandleMessage == " + string(msg))
		// TODO: 這邊要使用 event 來判斷是訊息還是下注

		// TODO: 若是下注則需要寫入 redis 做暫存

		// TODO: 發送訊息邏輯
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
	// 定義一個字串通道(開獎數字使用)
	message1 := make(chan string)
	// 定義一個字串通道(注碼計算使用)
	message1Return := make(chan string)

	// 負責開出數字的邏輯 Goroutine
	go func() {
		for {
			randServiceHandler := rand_generator(37) // 0 ~ 36
			result := <-randServiceHandler
			message1 <- strconv.Itoa(result)
			time.Sleep(time.Second * 20) // TODO: 目前設定為每 20 秒跑一次
		}
	}()

	// 負責接收通道發送出來的開獎數字訊息, 送出給前端網頁 Goroutine
	go func() {
		for result := range message1 {
			if result == "" {
				break
			} else {
				webSocket.Broadcast(NewMessage("roulette", "", result).GetByteMessage())
				fmt.Println(" -- this is a message : " + result)
				message1Return <- result
			}
		}
	}()

	//================================
	// 接收通道發送出的開獎數字, 進行注碼計算
	go func() {
		for item := range message1Return {
			fmt.Println("開獎後計算賠率以及返現 : " + item)
			// TODO: 抓取 redis 下注的人，計算賠率後回寫 redis 以及寫入資料庫
		}
	}()

	//================================
	// 啟動前先刪除整個 redis db
	deleteMsg, err2 := redisClient.FlushDB(context.Background()).Result()

	fmt.Println("===================")
	fmt.Println(deleteMsg) // 回傳是否刪除成功
	fmt.Println(err2)      // 是否有錯誤訊息
	fmt.Println("===================")

	// 處理 postgres sql 連線
	// Initialize connection string.  sslmode=disable , sslmode=require
	var connectionString string = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", HOST, USER, PASSWORD, DATABASE)

	// 初始化連線
	db, err := sql.Open("postgres", connectionString)
	checkError(err)

	err = db.Ping()
	checkError(err)
	fmt.Println("Successfully created connection to database")

	// 建立 repository
	fmt.Println("=================== Create repository Instance")
	rouletteRepo := _rouletteRepo.NewPostgresqlRouletteRepository(db)
	// 建立 usecase
	rouletteUsecase := _rouletteUsecase.NewRouletteUsecase(rouletteRepo)
	// 建立路由
	_rouletteHandlerHttpDelivery.NewRouletteHandler(server, rouletteUsecase)

	fmt.Println("http://localhost:8888")
	server.Run(":8888")
}
func checkError(err error) {
	if err != nil {
		panic(err)
	}
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
