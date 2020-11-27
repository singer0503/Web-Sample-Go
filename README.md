# 這是一個學習使用 Golang 開發 Web 應用的紀錄空間
- 01 原生 Golang 自帶的 http package 練習
- 02 顯示表單
- 03 如何避面 CROS
- 04 畫面和物件參數做連結 {{.Title}} {{.Content}}}
- 05 使用第三方 Framework（Gin）
- 06 使用 Gin 開發登入畫面
- 07 連接 MySQL 練習

傳統需要去設定 GOPATH
```go
GOPATH="/Users/Apple/Documents/Go/web-sample-go"
```
# 但是自從 Go Modules 的誕生

默認值，go命令會根據當前目錄来决定是否啟用modules功能。需要滿足兩種情形：
該專案目錄不在GOPATH/src/下
當前或上一層目錄存在 go.mod 檔案

```go
go mod init web-sample-go
```
執行之後可以看到會出現一個 go.mod 檔案