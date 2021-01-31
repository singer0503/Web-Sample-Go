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
	"log"
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

var CurrentTier = 100 // 目前設定，每一個籌碼的價值為 100 元

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

	// 測試用的
	server.GET("/test", handler.GetTest)
	// TODO: GetRouletteByBetID ...這個方法目前沒有使用到, 是一個 mvc 範例提供參考
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
			//獲取所有下注資料的 Redis hash 返回 map
			hashGetAll, _ := handler.RedisClient.HGetAll(context.Background(), BETS_HASH_KEY).Result()
			fmt.Println("HGetAll", hashGetAll)
			var msg = []byte("")

			//查看本局是否有人下注, 若有人下注才需要計算
			if len(hashGetAll) > 0 {
				handler.WebSocket.BroadcastFilter(msg, func(sessions *melody.Session) bool { // 線上有多少人就會跑幾次
					userSessionID, _ := sessions.Get(KEY)
					fmt.Println("====== s.Get(KEY)", userSessionID)
					// 確認該 user 是否有下注紀錄
					stringBet := hashGetAll[fmt.Sprintf("%v", userSessionID)]
					if stringBet != "" {
						fmt.Println("====我抓到這傢伙下注囉～", userSessionID)
						fmt.Println("====他的注碼～", stringBet)
						// 把字串轉換為 int陣列 Convert string to array of integers in golang
						var bets []int
						err := json.Unmarshal([]byte("["+stringBet+"]"), &bets)
						if err != nil {
							log.Fatal(err)
						}
						fmt.Printf("%v", bets)

						//檢查下注總數
						var bet = 0
						for i := 0; i < len(bets); i++ {
							if bets[i] != 0 {
								bet += bets[i]
							}

						}
						bet *= CurrentTier
						fmt.Println("===bet ", bet)

						//// TODO:檢查餘額

						// 計算勝負
						var win int = 0
						result, _ := strconv.Atoi(item)
						//如果壓中那個號碼是 36 倍的賠率
						if bets[result] != 0 {
							win += bets[result] * 36
						}
						//從 bets 37 以上開始算的意思是，0~36 是獨立的數字以在上面的邏輯就以算好，這邊迴圈是計算組合型投注的賠率！
						for i := 37; i < len(bets); i++ {
							if bets[i] != 0 {
								fmt.Println("sectormultipliers[i-37][result] === ", sectormultipliers[i-37][result])
								win += bets[i] * sectormultipliers[i-37][result] //計算陪率
							}
						}
						win *= 100 // 計算籌碼價值，目前都是 100 元
						win -= bet // 減掉投注額，就是贏回來的錢

						fmt.Println("下注(bet): ", bet, " 正負(win): ", win)
						returnBet := "[" + strconv.Itoa(bet) + "," + strconv.Itoa(win) + "]"
						betMsg := NewMessage("betUpdate", "", returnBet).GetByteMessage()
						//var oneSession melody.Session
						var newSessions = make([]*melody.Session, 0)
						newSessions = append(newSessions, sessions)
						//strUserSessionID := fmt.Sprintf("%v", userSessionID)
						//oneSession, _ := sessions.Get(strUserSessionID)
						fmt.Println("計算跑了幾次", userSessionID)
						handler.WebSocket.BroadcastMultiple(betMsg, newSessions)

						//一率不回訊息
						return false
					}
					// 沒有下注記錄則不需要發送訊息
					return false
				})
			}
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
	if message.Event == "bets" {
		fmt.Println(message.Content)
		fmt.Println("is bets message == " + fmt.Sprintf("%#v", message)) //直接把物件轉型成 string 並且列印
		// 印出是哪個人下注
		userSessionID, _ := s.Get(KEY)
		fmt.Println("====== s.Get(KEY)", userSessionID)
		// 將本局的下注資料儲存於 Redis
		d.RedisClient.HSet(context.Background(), BETS_HASH_KEY, userSessionID, message.Content)

	} else if message.Event == "message" {
		// TODO: 發送訊息邏輯
		listString, _ := d.GetRoom1List()
		d.WebSocket.BroadcastFilter(msg, func(sessions *melody.Session) bool {
			compareID, _ := sessions.Get(KEY)
			for _, value := range listString {
				if value == compareID {
					return true
				}
			}
			return false
		})
	}

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
	KEY           = "chat_id"
	WAIT          = "wait"
	ROOM1         = "room1"
	BETS_HASH_KEY = "room1-bets"
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

// 每個組合型注碼位置的賠率, 二維陣列表
var sectormultipliers = [12][37]int{
	{0, 0, 0, 3, 0, 0, 3, 0, 0, 3, 0, 0, 3, 0, 0, 3, 0, 0, 3, 0, 0, 3, 0, 0, 3, 0, 0, 3, 0, 0, 3, 0, 0, 3, 0, 0, 3}, //3rd column
	{0, 0, 3, 0, 0, 3, 0, 0, 3, 0, 0, 3, 0, 0, 3, 0, 0, 3, 0, 0, 3, 0, 0, 3, 0, 0, 3, 0, 0, 3, 0, 0, 3, 0, 0, 3, 0}, //2nd column
	{0, 3, 0, 0, 3, 0, 0, 3, 0, 0, 3, 0, 0, 3, 0, 0, 3, 0, 0, 3, 0, 0, 3, 0, 0, 3, 0, 0, 3, 0, 0, 3, 0, 0, 3, 0, 0}, //1st column
	{0, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, //1st 12
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, //2nd 12
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3}, //3rd 12
	{0, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, //1 to 18
	{0, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2}, //even
	{0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 0, 2, 0, 2, 0, 2, 0, 2, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 0, 2, 0, 2, 0, 2, 0, 2}, //Red
	{0, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 2, 0, 2, 0, 2, 0, 2, 0, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 2, 0, 2, 0, 2, 0, 2, 0}, //Black
	{0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 0}, //odd
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2}, //19 to 36
}
