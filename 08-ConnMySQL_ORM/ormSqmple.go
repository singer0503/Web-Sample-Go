// ORM 為 Object Relation Mapping 的縮寫，翻譯過來就是物件關聯對映。
// 參考網站：Day8 | 使用 ORM 與資料庫進行互動
// https://ithelp.ithome.com.tw/articles/10234820

package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	USERNAME = "root"
	PASSWORD = "1qaz!QAZ"
	NETWORK  = "tcp"
	SERVER   = "127.0.0.1"
	PORT     = 3306
	DATABASE = "demo"
)

type User struct {
	ID       int64  `json:"id" gorm:"primary_key;auto_increase'"`
	Username string `json:"username"`
	Password string `json:""`
}

func CreateUser(db *gorm.DB, user *User) error {
	return db.Create(user).Error
}

func FindUser(db *gorm.DB, id int64) (*User, error) {
	user := new(User)
	user.ID = id
	err := db.First(&user).Error
	return user, err
}

func main() {
	// 使用 orm 德方式連接資料庫
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("使用 gorm 連線 DB 發生錯誤，原因為 " + err.Error())
	}
	fmt.Println("conn test ok")

	// 使用 gorm 內建的 AutoMigrate 方法進行 Table 的建立，沒有 Table 時新增，Table 有更動時自動更新
	if err := db.AutoMigrate(new(User)); err != nil {
		panic("資料庫 Migrate 失敗，原因為 " + err.Error())
	}

	// 建立物件，使用該物件進行資料的新增
	user := &User{
		Username: "test2",
		Password: "test3",
	}
	if err := CreateUser(db, user); err != nil {
		panic("資料庫 Migrate 失敗原因為" + err.Error())
	}

	if user, err := FindUser(db, 4); err == nil {
		log.Print("查詢到 user 為", user)
	} else {
		panic("查詢 user 失敗，原因為" + err.Error())
	}

	// 更多語法可以參考 gorm 的官方文件。
	// https://gorm.io/
}
