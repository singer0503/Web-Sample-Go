// Day9 | 輕鬆管理程式的設定檔
// https://ithelp.ithome.com.tw/articles/10235053
package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func main() {
	fmt.Println("test")
	// 設定檔名
	viper.SetConfigName("app")
	// 設定型態
	viper.SetConfigType("yaml")
	// 設定路徑
	viper.AddConfigPath("./config")
	viper.SetDefault("application.port", 8080)
	// 讀取設定檔案
	err := viper.ReadInConfig()
	if err != nil {
		panic("讀取設定檔出現錯誤，原因為：" + err.Error())
	}
	// 抓出設定得值
	fmt.Println("application port = " + viper.GetString("application.port"))
	fmt.Println("application timeout read = " + viper.GetString("application.timeout.read"))
	fmt.Println("application apiBaseRoute = " + viper.GetString("application.apiBaseRoute"))
}
