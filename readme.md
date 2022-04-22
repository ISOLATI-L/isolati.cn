# ISOLATI
[ISOLATI](https://isolati.cn "")<br>
![](https://isolati.cn/files/1.jpg)
## 数据库配置
### 数据库信息配置文件
根目录下的SQL.config.ini设置数据库信息：
```ini
[SQL_Config]
server   = <地址>
port     = <端口>
user     = <用户名>
password = <密码>
database = <数据库名称>
```
### 建表语句
#### videos表
```sql
CREATE TABLE `videos` (
  `Vid` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `Vtitle` varchar(100) NOT NULL,
  `Vcontent` text,
  `Vcover` varchar(100) NOT NULL,
  `Vtime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`Vid`)
)
```
#### paragraphs表
```sql
CREATE TABLE `paragraphs` (
  `Pid` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `Ptitle` varchar(100) NOT NULL,
  `Pcontent` varchar(100) NOT NULL,
  `Ptime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`Pid`)
)
```
#### robots表
```sql
CREATE TABLE `robots` (
  `Rid` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `RuserAgent` text NOT NULL,
  `Rtime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`Rid`)
)
```
#### users表（暂未使用）
```sql
CREATE TABLE `users` (
  `Uid` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `Uname` varchar(50) NOT NULL,
  `Upassword` varchar(16) NOT NULL,
  `Uadmin` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`Uid`)
)
```
#### admins表
```sql
CREATE TABLE `admins` (
  `md5password` char(40) NOT NULL,
  PRIMARY KEY (`md5password`)
)
```
#### sessions表
```sql
CREATE TABLE `sessions` (
  `Sid` char(32) NOT NULL,
  `SlastAccessedTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `SmaxAge` int(10) unsigned NOT NULL DEFAULT '1800',
  `Sdata` json NOT NULL,
  PRIMARY KEY (`Sid`)
)
```
#### auto_delete_session事件
```sql
CREATE EVENT `auto_delete_session` ON SCHEDULE
EVERY 1 MINUTE DO
DELETE FROM
  sessions
WHERE
  (
    unix_timestamp(CURRENT_TIMESTAMP) - unix_timestamp(SlastAccessedTime)
  ) > SmaxAge
```
