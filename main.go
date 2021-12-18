package main

import (
	"fmt"
	"go-zentao-task/pkg/config"
	"go-zentao-task/pkg/db"
	"go-zentao-task/pkg/logging"
	"go-zentao-task/pkg/util"
	"go-zentao-task/pkg/zentaouser"
	"go-zentao-task/service/notification"
	"go-zentao-task/service/zentao"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var service = zentao.InitializeService()

type Task struct {
}

func setup(env string) {
	config.Setup(env)
	logging.Setup(env, logging.Stdout)
	db.Setup()
	service.User = zentaouser.Setup() //用户配置
	//rbac.Setup()
}

func getUsersByCsv() map[string]string {
	c1 := make(chan interface{})     //创建通道c1
	c2 := make(chan interface{})     //创建通道c2
	users := make(map[string]string) // 必须初始化 map, nil map 不能用来存放键值对
	go util.ReadCsv("xxx.csv", c1)   //启动一个goroutine, 读取csv
	go util.ReadCsv("xxx.csv", c2)   //启动另一个goroutine, 读取相同的csv
	//for i := range c1 {            // 从通道 c 中接收,遍历 通道c
	//	n := i.([]string) // 读取csv的 每一行都是[]string, 故可以 转换成 []string,
	//	users = append(users, n[0])
	//}
	// 用 select 实现多线程
	for len(users) < 4 { //循环直到切片有四个元素
		select {
		// 接收通道 c1 的结果
		case r := <-c1:
			if r != nil {
				n := r.([]string) // 读取csv的 每一行都是[]string, 故可以 转换成 []string,
				users[n[0]] = n[0]
			}
		// 接收通道 c2 的结果
		case r := <-c2:
			if r != nil {
				n := r.([]string)
				users[n[0]] = n[0]
			}
		default:
			fmt.Println("没获取到值")
		}
	}
	delete(users, "user") //去掉表头
	return users
}

func main() {
	uses := getUsersByCsv()
	fmt.Println(uses)
	var env = "development"
	setup(env) //初始化配置
	for _, u := range uses {
		service.User.Account = u
		service.UserLogin()
		time.Sleep(time.Duration(2) * time.Second)  //休眠2s
		es, _ := service.GetEstimateToday()         // 已用工时
		count, ids := service.ConsumeRecord(8 - es) //记录工时
		fmt.Println(count, ids)
		no := notification.NewNotification()                                                                // 创建报警实体
		no.SendNotification(service.User.Account, fmt.Sprintf("%.2f", es), fmt.Sprintf("%.2f", count), ids) // 发送报警，工时转成字符串
	}
	//监听终端quit命令
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("关闭服务 ...")
	log.Println("服务已退出")
}
