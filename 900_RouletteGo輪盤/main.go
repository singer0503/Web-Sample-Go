// Roulette 輪盤實作

// [Redis]-常用語法速查表
// https://www.dotblogs.com.tw/colinlin/2017/06/26/180604
package main

import (
	_rouletteHandlerHttpDelivery "RouletteGo/roulette/delivery/http"
	_rouletteRepo "RouletteGo/roulette/repository/postgresql"
	_rouletteUsecase "RouletteGo/roulette/usecase"
	"context"
	sql "database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	redis "github.com/go-redis/redis/v8"
	melody "gopkg.in/olahol/melody.v1"
)

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

	server.GET("/", func(c *gin.Context) {
		// 在http包使用的時候，註冊了/這個根路徑的模式處理，瀏覽器會自動的請求 favicon.ico，要注意處理，否則會出現兩次請求
		if c.Request.RequestURI == "/favicon.ico" {
			return
		}
		c.HTML(http.StatusOK, "index.html", nil)
	})

	webSocket := melody.New()

	//================================
	// 啟動前先刪除整個 redis db
	deleteMsg, err2 := redisClient.FlushDB(context.Background()).Result()

	fmt.Println("===================")
	fmt.Println(deleteMsg) // 回傳是否刪除成功
	fmt.Println(err2)      // 是否有錯誤訊息
	fmt.Println("===================")

	//初始化 viper，設定預設讀取環境變數
	// 設定檔名
	viper.SetConfigName("database")
	// 設定型態
	viper.SetConfigType("yaml")
	// 設定路徑
	viper.AddConfigPath("./config")
	// 讀取設定檔案
	err := viper.ReadInConfig()
	if err != nil {
		panic("讀取設定檔出現錯誤，原因為：" + err.Error())
	}
	// 處理 postgres sql 連線
	// Initialize connection string.  sslmode=disable , sslmode=require
	var connectionString string = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("postgreSQL.host"),
		viper.GetString("postgreSQL.user"),
		viper.GetString("postgreSQL.password"),
		viper.GetString("postgreSQL.database"))

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

	// 建立路由, gin , melody, redis 是從外部丟進去的
	_rouletteHandlerHttpDelivery.NewRouletteHandler(server, webSocket, redisClient, rouletteUsecase)

	fmt.Println("http://localhost:8888")
	server.Run(":8888")
}
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
