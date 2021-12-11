package main

import (
	"fmt"
	"go-zentao-task/model/zendao"
	"go-zentao-task/pkg/config"
	"go-zentao-task/pkg/db"
	"go-zentao-task/pkg/logging"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Task struct {
}

func setup(env string) {
	config.Setup(env)
	logging.Setup(env, logging.Stdout)
	//gredis.Setup()
	db.Setup()
	//session.Setup()
	//rbac.Setup()
}

func main() {
	fmt.Println("test")
	var env = "development"
	setup(env)
	var t Task
	t.View()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
	log.Println("Server exited")
}

func (*Task) View() {
	//获取任务
	task := zendao.NewTask()
	t, err := task.FindOne(10)
	if err != nil {
		fmt.Println( err.Error())
		return
	}
	fmt.Println(t)
	if t.ID == 0 {
		fmt.Println("任务不存在")
		return
	}
}
