package main

import (
	"FrontEndOJudger/caroline"
	"FrontEndOJudger/models"
	"FrontEndOJudger/pkg/setting"
	"FrontEndOJudger/pkg/ws"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"
)

func init() {
	// exec main judge logic
	setting.Setup()
	setting.Check()
	models.Setup()
	go ws.Setup()
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {

	// debug mode
	go func() {
		http.ListenAndServe("0.0.0.0:8899", nil)
	}()

	// start static file when run test chamber local
	go func() {
		if setting.JudgerSetting.TestChamberSwitch {
			http.Handle("/", http.FileServer(http.Dir(fmt.Sprintf("%s", setting.JudgerSetting.TestChamberBaseDir))))
			go http.ListenAndServe(fmt.Sprintf(":%s", setting.JudgerSetting.TestChamberPort), nil)
		}
	}()

	// start directly
	go func() {
		http.HandleFunc("/httpjudger", caroline.HttpJudger)
		http.HandleFunc("/screenshot", caroline.ScreenShot)
		go http.ListenAndServe(fmt.Sprintf(":%s", setting.JudgerSetting.HttpJudgerPort), nil)
	}()

	// fix expired judging submits while judger was crashed unexpectedly
	go func() {
		for {
			caroline.FixExpiredJudgingSubmits()
			time.Sleep(time.Duration(setting.JudgerSetting.SleepTime * 10) * time.Millisecond)
		}
	}()

	log.Printf("[SUCCESS] Project Caroline Judger Started 🎂")

	// main judge process
	go caroline.Judge()
	select{}
}
