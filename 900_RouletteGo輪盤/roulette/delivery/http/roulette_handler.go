package http

import (
	"RouletteGo/domain"
	swagger "RouletteGo/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gopkg.in/olahol/melody.v1"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// RouletteHandler ...
type RouletteHandler struct {
	RouletteUsecase domain.RouletteUsecase
	WebSocket       *melody.Melody
	RedisClient     *redis.Client
}

// NewRouletteHandler ... 路由控制
func NewRouletteHandler(server *gin.Engine, webSocket *melody.Melody, redisClient *redis.Client, rouletteUsecase domain.RouletteUsecase) {
	handler := &RouletteHandler{
		RouletteUsecase: rouletteUsecase,
		WebSocket:       webSocket,
		RedisClient:     redisClient,
	}

	server.GET("/ws", handler.GetWS)
	// 第一次使用 webSocket 連上來的處理
	webSocket.HandleConnect(handler.FirstConnect)
	// 收到訊息的處理
	webSocket.HandleMessage(handler.ReceiveMessage)
	// 關閉 webSocket 處理
	webSocket.HandleClose(handler.CloseConnect)

	server.GET("/test", handler.GetTest)

	// TODO: GetRouletteByBetID ...這個方法目前沒有使用到 mvc 範例
	server.GET("/api/v1/Roulette/:rouletteID", handler.GetRouletteByBetID)

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
}

func (d *RouletteHandler) GetWS(c *gin.Context) {
	d.WebSocket.HandleRequest(c.Writer, c.Request) // 外部丟入後在自身 handler
}

func (d *RouletteHandler) FirstConnect(session *melody.Session) {
	id := InitSession(session) // 產生 session key
	d.AddToRoom1List(id)       // 加入 redis room1 list 裡面
	// 這段只是測試用
	listString, _ := d.GetRoom1List()
	fmt.Println("listString===========", listString)
	msg := NewMessage("other", "對方已經", "加入輪盤局").GetByteMessage()

	d.WebSocket.BroadcastFilter(msg, func(session *melody.Session) bool { // 線上有多少人就會跑幾次
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
}

func (d *RouletteHandler) ReceiveMessage(s *melody.Session, msg []byte) {
	// TODO:
	var message Message
	json.Unmarshal(msg, &message)
	fmt.Println("HandleMessage == " + string(msg))
	// TODO: 這邊要使用 event 來判斷是訊息還是下注

	// TODO: 若是下注則需要寫入 redis 做暫存

	// TODO: 發送訊息邏輯
	listString, _ := d.GetRoom1List()
	d.WebSocket.BroadcastFilter(msg, func(session *melody.Session) bool {
		compareID, _ := session.Get(KEY)
		for _, value := range listString {
			if value == compareID {
				return true
			}
		}
		return false
	})
}

func (d *RouletteHandler) CloseConnect(session *melody.Session, i int, s string) error {
	id := GetSessionID(session)

	msg := NewMessage("other", "對方已經", "離開聊天室").GetByteMessage()
	d.RemoveRoom1ID(id)
	listString, _ := d.GetRoom1List()
	return d.WebSocket.BroadcastFilter(msg, func(session *melody.Session) bool {
		compareID, _ := session.Get(KEY)
		for _, value := range listString {
			if value == compareID {
				return true
			}
		}
		return false
	})
}

func (d *RouletteHandler) GetWebSocketHandleConnect(webSocket *melody.Melody) {

}

func (d *RouletteHandler) GetTest(c *gin.Context) {
	result := "{'msg':'test ok!'}"
	c.JSON(http.StatusOK, result)
}

// TODO: GetRouletteByBetID ...這個方法目前沒有使用到
func (d *RouletteHandler) GetRouletteByBetID(c *gin.Context) {
	rouletteID := c.Param("")

	anDigimon, err := d.RouletteUsecase.GetByID(c, rouletteID)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, &swagger.ModelError{
			Code:    3000,
			Message: "Internal error. Query digimon error",
		})
		return
	}

	c.JSON(200, &swagger.RoulletteInfo{
		Id:   anDigimon.ID,
		Name: anDigimon.Name,
	})
}

const (
	KEY   = "chat_id"
	WAIT  = "wait"
	ROOM1 = "room1"
)

type Message struct {
	Event   string `json:"event"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

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

func (d *RouletteHandler) GetRoom1List() ([]string, error) {
	return d.RedisClient.LRange(context.Background(), ROOM1, 0, -1).Result()
}
func (d *RouletteHandler) AddToRoom1List(id string) error {
	return d.RedisClient.LPush(context.Background(), ROOM1, id).Err()
}

func (d *RouletteHandler) RemoveRoom1ID(id string) {
	d.RedisClient.LRem(context.Background(), ROOM1, 0, id)
}

func InitSession(s *melody.Session) string {
	id := uuid.New().String()
	s.Set(KEY, id)
	return id
}
func GetSessionID(s *melody.Session) string {
	if id, isExist := s.Get(KEY); isExist {
		return id.(string)
	}
	return InitSession(s)
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
