# ISOLATI
[ISOLATI](https://isolati.cn "My Website")
## 数据库配置：
### 根目录下的SQL.config.ini存放数据库信息，格式如下：
```
[SQL_Config]
server   = <数据库IP地址>
port     = <端口>
user     = <用户名>
password = <密码>
database = <数据库名称>
```
### 目前数据库有3张表：
#### videos表：
```
CREATE TABLE `videos` (
  `Vid` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `Vtitle` varchar(100) NOT NULL,
  `Vcontent` text,
  `Vcover` varchar(100) NOT NULL DEFAULT '',
  `Vtime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`Vid`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8
```
#### users表（暂未用到）：
```
CREATE TABLE `users` (
  `Uid` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `Uname` varchar(50) NOT NULL,
  `Upassword` varchar(16) NOT NULL,
  `Uadmin` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`Uid`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8
```
#### sessions表（暂未用到）：
```
CREATE TABLE `sessions` (
  `Sid` char(32) NOT NULL,
  `SlastAccessedTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `SmaxAge` int(10) unsigned NOT NULL DEFAULT '1800',
  `Sdata` json NOT NULL,
  PRIMARY KEY (`Sid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8
```
