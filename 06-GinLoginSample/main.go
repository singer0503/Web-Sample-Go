// 參考網站：Day6 | 透過 golang 實作一個簡單的登入功能
// https://ithelp.ithome.com.tw/articles/10234298
package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	server.LoadHTMLGlob("template/html/*")
	//設定靜態資源的讀取
	server.Static("/assets", "./template/assets")
	server.GET("/login", LoginPage)
	server.POST("/login", LoginAuth)
	server.Run(":8888")
}

// 這邊有遇到問題，就是同目錄夾下的 main 包，居然是不可見的，因為 go run main.go ，只會針對這個檔案做編譯
// 就會出現錯誤
// # command-line-arguments
// ./main.go:14:23: undefined: LoginPage
// ./main.go:15:24: undefined: LoginAuth
// 編譯成可執行檔 go build
// 在執行那個可執行檔即可
