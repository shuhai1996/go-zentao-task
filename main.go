package main

import (
	"fmt"
	"go-zentao-task/pkg/config"
	"go-zentao-task/pkg/db"
	"go-zentao-task/pkg/logging"
	"go-zentao-task/pkg/zentaouser"
	"go-zentao-task/service/zentao"
	"log"
	"os"
	"os/signal"
	"syscall"
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

func main() {
	var env = "development"
	setup(env) //初始化配置
	//es,_ := service.GetEstimateToday()// 已用工时
	//es := service.GetAllTaskNotDone()
	es := service.ConsumeRecord()
	fmt.Println(es)
	//url := util.GetRobotUrl()
	//no := notification.NewNotification()// 创建报警实体
	//no.SendNotification(url, service.User.Account, strconv.FormatFloat(float64(es), 'f', 10, 32))// 发送报警，工时转成字符串
	//监听终端quit命令
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
	log.Println("Server exited")
}
