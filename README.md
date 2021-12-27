## go-zentao-task

#### 配置
根目录下创建conf.ini文件，添加如下配置
```
[development]
db.host =localhost
db.port =3306
db.username =root
db.password =
db.database =
```
#### 包管理
删除go.mod和go.sum文件，执行下列命令：
```
go mod init go-zentao-task
go mod vendor //将下载到GOPATH的包复制到当前项目的vendor目录下
```
#### 运行本项目
执行如下命令即可运行
```
go build main.go
./main
```

> 参考链接
> 
> [https://gin-gonic.com/](https://gin-gonic.com/)
