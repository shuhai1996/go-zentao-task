package main

import (
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
	service.User = zentaouser.Setup()
	//session.Setup()
	//rbac.Setup()
}

func main() {
	var env = "development"
	setup(env)
	service.GetEstimateToday()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
	log.Println("Server exited")
}
