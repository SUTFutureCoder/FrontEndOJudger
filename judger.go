package main

import (
	"FrontEndOJudger/caroline"
	"FrontEndOJudger/models"
	"FrontEndOJudger/pkg/setting"
	"github.com/takama/daemon"
	"log"
	"os"
)

type Service struct {
	daemon.Daemon
}

const (
	name = "feonlinejudger"
	description = "FrontEndOnlineJudger"
)

func (service *Service) Manage() (string, error) {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		default:
			log.Fatalf("Usage: feonlinejudger install | remove | start | stop | status")
		}
	}

	// 执行核心judger逻辑
	setting.Setup()
	models.Setup()
	caroline.Judge()

	return "", nil
}

func main() {
	srv, err := daemon.New(name, description)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	service := &Service{srv}
	execStr, err := service.Manage()
	if err != nil {
		log.Fatalf("status: %v Error %v", execStr, err)
	}
	if execStr != "" {
		log.Println(execStr)
	}

}
