安裝
在 mac 環境可以使用 HomeBrew 來進行安裝，首先我們先將進行套件更新
`brew update`

接著進行 redis 的安裝
`brew install redis`

啟動
Redis 的啟動也可以使用 brew 指令進行
`brew services start redis`

測試
可以透過 ping 的方式測試看看 redis-server 是否存活

`redis-cli ping`
這時如果 command line 看到回傳 PONG，代表安裝成功拉！

有时候会有中文乱码。
要在 redis-cli 后面加上 --raw

`redis-cli --raw`
就可以避免中文乱码了。

`config set requirepass "admin"`