package main

// 參考網站：Day5 | Gin - 好用的 web framework
// https://ithelp.ithome.com.tw/articles/10234075

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type IndexData struct {
	Title   string
	Content string
}

func test(c *gin.Context) {
	data := new(IndexData)
	data.Title = "首頁"
	data.Content = "我的第一個首頁"
	c.HTML(http.StatusOK, "index.html", data)
}
func main() {
	server := gin.Default()           // 建立 gin 的 instance，可以把這個 instance 想像成是 server 的實例
	server.LoadHTMLGlob("template/*") // 這邊要先跟 gin 註冊好 template 的位置，好讓他去找
	server.GET("/", test)             // 設定 routing ，gin 原生提供了封裝各種 http method 的 routing
	server.Run(":8888")               // 啟動 gin server，只需要單純的把 gin 的 instance 進行 Run() 方法即可
}
