## go-zentao-task

### 运行本项目
#### 编译
执行如下命令即可完成编译
```
go build
```
#### 运行
执行如下命令即可运行

#### 包管理
新的项目，复制完代码以后，全局替换go-zentao-task为你的项目名称，删除go.mod和go.sum文件，执行下列命令：
```
go mod init go-zentao-task
go build //更新go.mod文件，包版本管理，下载远程包并编译
go mod vendor //将下载到GOPATH的包复制到当前项目的vendor目录下
go build -mod=vendor //依赖当前项目下vendor文件夹中的包进行编译，为了避免在线编译时包下载失败，都要基于此模式编译
```

> 参考链接 [https://juejin.im/post/5d8ee2db6fb9a04e0b0d9c8b](https://juejin.im/post/5d8ee2db6fb9a04e0b0d9c8b)
