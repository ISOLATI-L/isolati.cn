# ISOLATI
[ISOLATI](https://isolati.cn "My Website")<br>
![](https://isolati.cn/files/1.jpg)
## 数据库配置
### 数据库信息配置文件
根目录下的SQL.config.ini设置数据库信息：
```
[SQL_Config]
server   = <地址>
port     = <端口>
user     = <用户名>
password = <密码>
database = <数据库名称>
```
### 数据库设置
#### videos表
```
CREATE TABLE `videos` (
  `Vid` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `Vtitle` varchar(100) NOT NULL,
  `Vcontent` text,
  `Vcover` varchar(100) NOT NULL DEFAULT '',
  `Vtime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`Vid`)
)
```
#### robots表
```
CREATE TABLE `robots` (
  `Rid` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `RuserAgent` text NOT NULL,
  `Rtime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`Rid`)
)
```
#### users表（暂未使用）
```
CREATE TABLE `users` (
  `Uid` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `Uname` varchar(50) NOT NULL,
  `Upassword` varchar(16) NOT NULL,
  `Uadmin` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`Uid`)
)
```
#### admins表
```
CREATE TABLE `admins` (
  `md5password` char(40) NOT NULL,
  PRIMARY KEY (`md5password`)
)
```
#### sessions表（暂未使用）
```
CREATE TABLE `sessions` (
  `Sid` char(32) NOT NULL,
  `SlastAccessedTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `SmaxAge` int(10) unsigned NOT NULL DEFAULT '1800',
  `Sdata` json NOT NULL,
  PRIMARY KEY (`Sid`)
)
```
#### auto_delete_session事件（暂未使用）
```
CREATE EVENT `auto_delete_session` ON SCHEDULE
EVERY 1 MINUTE DO
DELETE FROM
  sessions
WHERE
  (
    unix_timestamp(CURRENT_TIMESTAMP) - unix_timestamp(SlastAccessedTime)
  ) > SmaxAge
```
