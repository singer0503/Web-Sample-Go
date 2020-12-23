package main

import (
	"fmt"
	rand "math/rand"
	"strconv"
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

func BatchGeneratorRand(rangeNumber int) chan int {
	rand.Seed(time.Now().UnixNano()) // 設定種子, 不然下次開啟也是一樣（利用當前時間的 UNIX 時間戳初始化 rand package）
	result := make(chan int)
	go func(x int) {
		for {
			result <- rand.Intn(x)
		}
	}(rangeNumber)
	return result
}

func main() {

	// 定义一个字符型的通道
	// message := make(chan string)

	// go func() {
	// 	for i := 1; i <= 5; i++ {
	// 		if i == 5 {
	// 			message <- ""
	// 		} else {
	// 			rand.Seed(time.Now().UnixNano()) // 設定種子, 不然下次開啟也是一樣（利用當前時間的 UNIX 時間戳初始化 rand package）
	// 			x := rand.Intn(100)
	// 			message <- ("this is a message : " + strconv.Itoa(x))
	// 		}
	// 		time.Sleep(time.Second * 5)
	// 	}
	// }()

	// // 接收通道发送的消息
	// for result := range message {
	// 	if result == "" {
	// 		break
	// 	} else {
	// 		fmt.Println(result)
	// 	}
	// }

	// 定義一個字串通道
	message1 := make(chan string)
	go func() {
		for {
			rand_service_handler := rand_generator(50) // 0 ~ 49
			result := <-rand_service_handler
			message1 <- strconv.Itoa(result)
			time.Sleep(time.Second * 2)
		}
	}()
	// // 定義一個字串通道
	// message2 := make(chan string)
	// go func() {
	// 	for {
	// 		rand_service_handler := rand_generator(50)
	// 		//fmt.Printf(" 1 -- this is a message : "+"%d\n", <-rand_service_handler)
	// 		result := <-rand_service_handler
	// 		message1 <- strconv.Itoa(result)
	// 		time.Sleep(time.Second * 2)
	// 	}
	// }()

	// go func() {
	// 	for {
	// 		rand_service_handler := rand_generator(50)
	// 		fmt.Printf(" 2 -- this is a message : "+"%d\n", <-rand_service_handler)
	// 		time.Sleep(time.Second * 2)
	// 	}
	// }()

	// 接收通道发送的消息
	for result := range message1 {
		if result == "" {
			break
		} else {
			fmt.Println(" 1 -- this is a message : " + result)
		}
	}
	//time.Sleep(time.Second * 100)

	// startdt := time.Now()
	// fmt.Println("Current date and time is: " + startdt.Format("01-02-2006 15:04:05.000000"))
	// // 生成随机数作为一个服务
	// for i := 0; i < 10; i++ {
	// 	rand_service_handler := rand_generator(100)
	// 	fmt.Printf("%d\n", <-rand_service_handler)
	// }
	// // 从服务中读取随机数并打印

	// enddt := time.Now()
	// fmt.Println("Current date and time is: " + enddt.Format("01-02-2006 15:04:05.000000"))
}
