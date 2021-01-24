package main

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
	//if _, err := db.Exec(sql); err != nil {
	//	fmt.Println("建立 Table 發生錯誤", err)
	//	return err
	//}
	//fmt.Println("建立 Table 成功！")
	//return nil
	retult, err := db.Exec(sql)

	if retult != nil {
		fmt.Println(retult.LastInsertId())
		fmt.Println(retult.RowsAffected())
	}

	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func main() {
	fmt.Println("test")
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

	CreateTable(db)

}
