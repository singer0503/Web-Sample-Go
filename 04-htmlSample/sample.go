package main

// 參考網站： Day4 | 無痛使用 Golang 打造屬於自己的網頁、Day5 | Gin - 好用的 web framework
// https://ithelp.ithome.com.tw/articles/10233981
// https://ithelp.ithome.com.tw/articles/10234075
import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type IndexData struct {
	Title   string
	Content string
}

func test(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("/index.html"))
	data := new(IndexData)
	data.Title = "首頁"
	data.Content = "我的第一個首頁"
	tmpl.Execute(w, data) // 會把 index.html 畫面和物件參數做連結 {{.Title}} {{.Content}}
}

func main() {
	http.HandleFunc("/", test)
	http.HandleFunc("/index", test)
	fmt.Println("start ok! http://localhost:8888")
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
