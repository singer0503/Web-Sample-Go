# 這是一個學習使用 Golang 開發 Web 應用的紀錄空間
- 01 原生 Golang 自帶的 http package 練習
- 02 顯示表單
- 03 如何避面 CROS
- 04 畫面和物件參數做連結 {{.Title}} {{.Content}}}
- 05 使用第三方 Framework（Gin）
- 06 使用 Gin 開發登入畫面
- 07 連接 MySQL 練習
- 08 連接 MySQL 並使用 gorm 
- 09 env 練習
- 10 i18n 練習
- 11 使用者管理 Web api 開發練習，並實作乾淨的架構
- 18_WebSocket
- 21_WebSocket+Redis一對一聊天室
- 23_UniTest
- 100_ginJWT 範例
- 101 Random 隨機數應用範例
- 102_css-roulette-wheel前端參考範例
- 103_bet_plate前端參考範例
- 900_RouletteGo 輪盤       ###輪盤作品請直接參考這個！！！ 

傳統需要去設定 GOPATH
```go
GOPATH="/Users/Apple/Documents/Go/web-sample-go"
```
# 但是自從 Go Modules 的誕生
首先要先設定 GO111MODULE 環境變數，總共可以三種不同的值：

`auto`
默認值，go命令會根據當前目錄来决定是否啟用modules功能。需要滿足兩種情形：
- 該專案目錄不在GOPATH/src/下
- 當前或上一層目錄存在go.mod檔案

`on`
go 命令會使用 modules，而不會GOPATH目錄下查找。

`off`
go 命令將不會支持 module 功能，尋找套件如以前 GOPATH 的做法去尋找。

反正設定為 `on` 就對了，`go mod` 就像是 java 的 gradle maven、dotnet 的 Nuget 一樣，原生搭載套件管理器！！

```go
go env -w GO111MODULE=on
```

MacOS 或者 Linux 下開啟 GO111MODULE 的命令為：

export GO111MODULE=on 或者 export GO111MODULE=auto

```go
go mod init web-sample-go
```
執行之後可以看到會出現一個 go.mod 檔案

假設現在我要引入GitHub上的gin-gonic/gin的套件，如下定義：

```go
module modtest

go 1.13

require github.com/gin-gonic/gin v1.5.0
```

再執行以下指令：
`go mod download`

會將需要的套件安裝在 GOPATH/pkg/mod 資料夾裡面。而且會發現出現一個 go.sum 的檔案，這個檔案基本上用來記錄套件版本的關係，確保是正確的，是不太需要理會的。

包(Package)的名稱 為包所在的資料夾名稱

# Redius
抬起頭來FLUSHALL可能太過猛烈了。FLUSHDB是僅刷新數據庫的一個。
FLUSHALL將清除整個服務器。就像服務器上的每個數據庫一樣。
由於問題是關於刷新數據庫的，所以我認為這是一個重要的區別，值得單獨回答。

清除整個 Redius Server
```
redius-cli FLUSHALL
```

清除 Redius 該資料庫
```
redius-cli FLUSHDB
```

查詢 全部的 key
```bash
323-Maxhuang:~ Apple$ redis-cli
127.0.0.1:6379> keys *
1) "61243c22-4351-42eb-aeb4-622bc6b28ac0"
2) "a1969a7a-e1d2-456c-ade8-4e16dd7bfae1"
127.0.0.1:6379> 

```

首先應該明白報這個錯誤說明了你用的jedis方法與redis服務器中存儲數據的類型存在衝突。

例如：數據庫中有一個key是usrInfo的數據存儲的是Hash類型的，但是你使用jedis執行數據

操作的時候卻使用了非Hash的操作方法，比如Sorted Sets裡的方法。此時就會報

ERR Operation against a key holding the wrong kind of value這個錯誤！

問題解決：

先執行一條如下命令，usrInfo為其中的一個key值。

redis 127.0.0.1:6379>type usrInfo

此時會顯示出該key存儲在現在redis服務器中的類型，例如：

redis 127.0.0.1:6379>hash
則表示key為usrInfo的數據是以hash類型存儲在redis服務器裡的，此時操作這個數據就必須使用hset、hget等操作方法。

例如 1 對 1 聊天室的 wait 是用 list
```
127.0.0.1:6379> keys *
1) "wait"
127.0.0.1:6379> get wait
(error) WRONGTYPE Operation against a key holding the wrong kind of value
127.0.0.1:6379> type wait
list
127.0.0.1:6379> LLEN wait   
(integer) 1
127.0.0.1:6379> LINDEX wait 0
"8f67935d-ba01-42a7-80b4-5e9a286dab4b"
127.0.0.1:6379> 
```
[Redis]-常用語法速查表 https://www.dotblogs.com.tw/colinlin/2017/06/26/180604

一次取得該 list 所有資料
```
LRANGE room1 0 -1
```


刪除 list 裡面的某一筆資料
https://redis.io/commands/lrem
```
LREM room1 0 "de72bfa2-e8fd-4174-9a8f-15d044925591"
```

### css 設計參考
https://wcc723.github.io/css/2017/07/21/css-flex/

### Less 介紹
https://www.oxxostudio.tw/articles/201601/css-less-01.html

### go 的 json.Unmarshal 可以把 json 字串轉成 struct，而 json.Marshal 可以將 struct 轉成 json 字串．
https://ithelp.ithome.com.tw/articles/10205062
https://www.flysnow.org/2018/11/05/golang-concat-strings-performance-analysis.html