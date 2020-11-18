安裝
golang 在使用外來的 package 時就使用 go get 的指令就可以輕鬆的 package 下載下來並且安裝
讓我們使用以下的 command 將 gin 安裝到自己的 package 當中

go get github.com/gin-gonic/gin
快速開始使用
匯入 gin package
一開始要使用 gin 不外乎與 net/http package 一樣需要進行 import，因此我們就先將這兩個 package 進行 import

import (
    "github.com/gin-gonic/gin"
    "net/http"
)
加入 request handler
如同簡介中所提到的， gin 在處理 request 與 response 都用一個萬用的 gin.Context 就可以拉，就讓我們延續昨天的範例，回傳一個 HTML

func test(c *gin.Context) {
	data := new(IndexData)
	data.Title = "首頁"
	data.Content = "我的第一支 gin 專案"
	c.HTML(http.StatusOK, "index.html", data)
}
