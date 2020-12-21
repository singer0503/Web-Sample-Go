package main

import (
	"fmt"
	rand "math/rand"
	"time"
)

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

func main() {

	// 產生十個 0 ~ 99 隨機數

	// rand.Seed(time.Now().UnixNano()) // 設定種子, 不然下次開啟也是一樣（利用當前時間的 UNIX 時間戳初始化 rand package）
	// for i := 0; i < 10; i++ {
	// 	x := rand.Intn(100)
	// 	fmt.Println(x)
	// }
	// enddt := time.Now()
	// fmt.Println("Current date and time is: " + enddt.Format("01-02-2006 15:04:05.000000"))

	startdt := time.Now()
	fmt.Println("Current date and time is: " + startdt.Format("01-02-2006 15:04:05.000000"))
	// 生成随机数作为一个服务
	for i := 0; i < 10; i++ {
		rand_service_handler := rand_generator(100)
		fmt.Printf("%d\n", <-rand_service_handler)
	}
	// 从服务中读取随机数并打印

	enddt := time.Now()
	fmt.Println("Current date and time is: " + enddt.Format("01-02-2006 15:04:05.000000"))
}
