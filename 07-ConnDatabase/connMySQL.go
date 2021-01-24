package main

// 參考網站：Day7 | 使用 GoLang 與資料庫進行互動
// https://ithelp.ithome.com.tw/articles/10234657
import (
	"database/sql"
	"fmt"

	_ "github.com/model-sql-driver/mysql" // 使用 mysql Driver, 使用 _ 可以去呼叫裡面的 init 方法進行初始化
)

const (
	USERNAME = "root"
	PASSWORD = "1qaz!QAZ"
	NETWORK  = "tcp"
	SERVER   = "127.0.0.1"
	PORT     = 3306
	DATABASE = "demo"
)

func CreateTable(db *sql.DB) error {
	sql := `CREATE TABLE IF NOT EXISTS users (
		id INT(4) PRIMARY KEY AUTO_INCREMENT NOT NULL,
		username VARCHAR(64),
		password VARCHAR(64)
	);`
	if _, err := db.Exec(sql); err != nil {
		fmt.Println("建立 Table 發生錯誤", err)
		return err
	}
	fmt.Println("建立 Table 成功！")
	return nil
}

func InsertUser(DB *sql.DB, username, password string) error {
	_, err := DB.Exec("insert INTO users(username,password) values(?,?)", username, password)
	if err != nil {
		fmt.Printf("建立使用者失敗，原因是：%v", err)
		return err
	}
	fmt.Println("建立使用者成功！")
	return nil
}

type User struct {
	ID       string
	Username string
	Password string
}

func QueryUser(db *sql.DB, username string) {
	user := new(User)
	row := db.QueryRow("select * from users where username = ?", username)
	if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		fmt.Printf("映射使用者失敗，原因為：%v\n", err)
		return
	}
	fmt.Println("查詢使用者成功", *user)
}

func main() {
	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	fmt.Println(conn)
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
	//CreateTable(db) // 建立資料表
	//InsertUser(db, "test", "test") // 新增資料
	QueryUser(db, "test") // 查詢資料
}
