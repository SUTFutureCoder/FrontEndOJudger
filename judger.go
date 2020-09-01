package main

import (
	"FrontEndOJudger/caroline"
	"FrontEndOJudger/models"
	"FrontEndOJudger/pkg/setting"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

func main() {

	// 执行核心judger逻辑
	setting.Setup()
	models.Setup()
	runtime.GOMAXPROCS(runtime.NumCPU())

	// start static file when run test chamber local
	if setting.JudgerSetting.TestChamberSwitch {
		http.Handle("/", http.FileServer(http.Dir(fmt.Sprintf("./%s", setting.JudgerSetting.TestChamberDir))))
		go http.ListenAndServe(fmt.Sprintf(":%s", setting.JudgerSetting.TestChamberPort), nil)
	}

	for {
		caroline.Judge()
		time.Sleep(time.Duration(setting.JudgerSetting.SleepTime) * time.Millisecond)
	}

}
