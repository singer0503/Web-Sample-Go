連接 MySQL

建立資料庫與使用者帳號，使用 MySQL 的 root 管理者帳號登入：

```
mysql -u root -p
```

在 MySQL/MariaDB 中新增資料庫：
## 新增資料庫
```sql
CREATE DATABASE `test`;
```

```sql
show databases;

use test;

CREATE TABLE `userinfo` (
    `uid` INT(10) NOT NULL AUTO_INCREMENT,
    `username` VARCHAR(64) NULL DEFAULT NULL,
    `department` VARCHAR(64) NULL DEFAULT NULL,
    `created` DATE NULL DEFAULT NULL,
    PRIMARY KEY (`uid`)
);

CREATE TABLE `userdetail` (
    `uid` INT(10) NOT NULL DEFAULT '0',
    `intro` TEXT NULL,
    `profile` TEXT NULL,
    PRIMARY KEY (`uid`)
);

```


這樣就會新增一個新的 my_db 資料庫。

## 新增使用者，設定密碼
新增一個 MySQL 資料庫使用者 my_user，並將密碼設定為 my_password：
```
CREATE USER 'my_user'@'localhost' IDENTIFIED BY 'my_password';
```

授予 my_user 帳號在 my_db 資料庫上面的所有權限，也就是讓 my_user 可以對 my_db 資料庫進行任何操作：
設定使用者權限
```
GRANT ALL PRIVILEGES ON my_db.* TO 'my_user'@'localhost';
```

使用 GRANT 設定好帳號的權限之後，馬上就會生效。接著就可以離開 MySQL 資料庫，重新以新的帳號登入使用了：
```
mysql -u my_user -p
```
