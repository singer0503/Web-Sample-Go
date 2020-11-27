連接 MySQL

```bash
mysql -u root -p
1qaz!QAZ

```

# **建立使用者與資料庫**

使用 `mysql` 指令進入控制台後，建立一個名為 `demo` 的資料庫

```
CREATE DATABASE demo CHARACTER SET utf8 COLLATE utf8_general_ci;

```

建立帳號為 `demo`，密碼為 `demo123`

```
GRANT ALL PRIVILEGES ON demo.* TO 'demo1'@'%' IDENTIFIED BY 'demo123'  WITH GRANT OPTION;
FLUSH PRIVILEGES;

```

# **操作 mysql**

要透過程式語言操作資料庫，最常見的方法就是使用 `driver`，golang 原生有提供關於 sql 的抽象介面 [database/sql](https://golang.org/pkg/database/sql/)，後來有人利用他封裝了 `mysql` 的 driver - [go-sql-driver](https://github.com/go-sql-driver/mysql)，接下來我們會利用這個 package 進行練習。