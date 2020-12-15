// 透過 Go 與 Redis 進行互動
// https://ithelp.ithome.com.tw/articles/10240302
package main

import (
	"context"
	"fmt"
	"log"

	redis "github.com/go-redis/redis/v8"
)

func main() {
	fmt.Println("test")

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		//Password: "a12345", // no password set
		//DB:       0,        // use default DB
	})

	pong, err := rdb.Ping(context.Background()).Result()
	if err == nil {
		log.Println("redis 回應成功，", pong)
	} else {
		log.Fatal("redis 無法連線，錯誤為", err)
	}
}
