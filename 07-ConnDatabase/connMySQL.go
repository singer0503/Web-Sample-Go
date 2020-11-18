package main

// 參考網站：Day7 | 使用 GoLang 與資料庫進行互動
// https://ithelp.ithome.com.tw/articles/10234657
import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // 使用 mysql Driver, 使用 _ 可以去呼叫裡面的 init 方法進行初始化
)

const (
	USERNAME = "root"
	PASSWORD = "1qaz!QAZ"
	NETWORK  = "tcp"
	SERVER   = "127.0.0.1"
	PORT     = 3306
	DATABASE = "demo"
)

func main() {
	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	db, err := sql.Open("mysql", conn)
	if err != nil {
		fmt.Println("開啟 MySQL 連線發生錯誤，原因為：", err)
		return
	}
	if err := db.Ping(); err != nil {
		fmt.Println("資料庫連線錯誤，原因為：", err.Error())
		return
	}
	defer db.Close()
}
