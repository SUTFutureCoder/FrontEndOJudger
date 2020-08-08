package main

import (
	"FrontEndOJudger/caroline"
	"FrontEndOJudger/models"
	"FrontEndOJudger/pkg/setting"
	"runtime"
	"time"
)

func main() {

	// 执行核心judger逻辑
	setting.Setup()
	models.Setup()
	runtime.GOMAXPROCS(runtime.NumCPU())

	for {
		caroline.Judge()
		time.Sleep(time.Duration(setting.JudgerSetting.SleepTime) * time.Millisecond)
	}

}
