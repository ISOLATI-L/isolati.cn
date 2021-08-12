# ISOLATI
[ISOLATI](https://isolati.cn "My Website")<br/>
![](https://isolati.cn/files/1.jpg)
## 数据库配置
### 数据库信息配置文件
根目录下的SQL.config.ini存放数据库信息，格式如下：
```
[SQL_Config]
server   = <数据库IP地址>
port     = <端口>
user     = <用户名>
password = <密码>
database = <数据库名称>
```
### 数据库表
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
