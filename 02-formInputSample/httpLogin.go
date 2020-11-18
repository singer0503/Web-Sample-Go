package main

// 處理表單輸入
// https://willh.gitbook.io/build-web-application-with-golang-zhtw/04.0/04.1
import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var folderPath string

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析 url 傳遞的參數，對於 POST 則解析 HTTP 回應內容的主體（request body）
	//注意 : 如果沒有呼叫 ParseForm 方法，下面無法取得表單的資料
	fmt.Println(r.Form) //這些資訊是輸出到伺服器端的列印資訊
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie22!") //這個寫入到 w 的是輸出到客戶端的
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //取得請求的方法
	if r.Method == "GET" {
		fmt.Println("get")
		//t, err := template.ParseFiles(folderPath + "/login.gtpl") // 加上檔案所在位置
		t, err := template.ParseFiles("/login.gtpl")
		if err != nil {
			fmt.Println("===== template: ", err)
		}
		log.Println(t.Execute(w, nil))
	} else {
		r.ParseForm() // 要拿 body 的資料，一定要這行 讓他解析
		//請求的是登入資料，那麼執行登入的邏輯判斷
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
	}
}

func main() {
	http.HandleFunc("/", sayhelloName) //設定存取的路由
	http.HandleFunc("/login", login)   //設定存取的路由
	fmt.Println("http://localhost:9091")
	fmt.Println("http://localhost:9091/login")

	// 解析是哪一種作業系統 mac, windows, linux
	ex, errOS := os.Executable()
	if errOS != nil {
		panic(errOS)
	}
	// 根據不同系統, 取得當前檔案的所在位置
	folderPath = filepath.Dir(ex)
	fmt.Println(folderPath)

	err := http.ListenAndServe(":9091", nil) //設定監聽的埠
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
