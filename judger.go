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

	// exec main judge logic
	setting.Setup()
	setting.Check()
	models.Setup()
	runtime.GOMAXPROCS(runtime.NumCPU())

	// start static file when run test chamber local
	if setting.JudgerSetting.TestChamberSwitch {
		http.Handle("/", http.FileServer(http.Dir(fmt.Sprintf("%s/%s", setting.JudgerSetting.TestChamberBaseDir, setting.JudgerSetting.TestChamberDir))))
		go http.ListenAndServe(fmt.Sprintf(":%s", setting.JudgerSetting.TestChamberPort), nil)
	}

	// fix expired judging submits while judger was crashed unexpectedly
	go func() {
		for {
			caroline.FixExpiredJudgingSubmits()
			time.Sleep(time.Duration(setting.JudgerSetting.SleepTime * 10) * time.Millisecond)
		}
	}()

	for {
		// main judge process
		caroline.Judge()
		time.Sleep(time.Duration(setting.JudgerSetting.SleepTime) * time.Millisecond)
	}

}
