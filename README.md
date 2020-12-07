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
- 12 WebSocket

傳統需要去設定 GOPATH
```go
GOPATH="/Users/Apple/Documents/Go/web-sample-go"
```
# 但是自從 Go Modules 的誕生
首先要先設定GO111MODULE環境變數，總共可以三種不同的值：

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