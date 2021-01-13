package setting

import (
	"fmt"
	"github.com/go-ini/ini"
	"log"
	"os"
)

type Judger struct {
	JudgerSum          int
	SleepTime          int

	TestChamberSwitch bool
	TestChamberBaseDir string
	TestChamberDir string
	TestChamberAddr string
	TestChamberPort string

}

var JudgerSetting = &Judger{}

type Database struct {
	Type         string
	User         string
	Password     string
	Host         string
	Name         string
	TablePrefix  string
	MaxOpenConns int
	MaxIdleConns int
}

var DatabaseSetting = &Database{}

var cfg *ini.File

func Setup() {
	var err error
	cfg, err = ini.Load("conf/judger.ini")
	if err != nil {
		log.Fatalf("Judger setup failed, place check [conf/judger.ini] file exist error:%v", err)
	}

	mapTo("database", DatabaseSetting)
	mapTo("judger", JudgerSetting)
}

func mapTo(s string, v interface{}) {
	err := cfg.Section(s).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", s, err)
	}
}

/**
 check and repair settings if not correct
 */
func Check() {
	checkJudger()
}

func checkJudger() {
	// check judger TestChamberBaseDir
	testchamber := fmt.Sprintf("%s/%s", JudgerSetting.TestChamberBaseDir, JudgerSetting.TestChamberDir)
	_, err := os.Stat(testchamber)
	if err == nil || os.IsExist(err) {
		return
	}
	// try create dir if not exist
	err = os.MkdirAll(testchamber, 0777)
	if err == nil {
		return
	}

	// failed to create dir but try user home dir
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Can not create target dir and your home dir, please check your TestChamberBaseDir config(%s) in conf/judger.ini was set correctly.", JudgerSetting.TestChamberBaseDir)
		return
	}

	// correct to home dir
	JudgerSetting.TestChamberBaseDir = fmt.Sprintf("%s/FrontEndOJudger", homeDir)
	testchamber = fmt.Sprintf("%s/%s", JudgerSetting.TestChamberBaseDir, JudgerSetting.TestChamberDir)
	err = os.MkdirAll(testchamber, 0777)
	if err != nil {
		log.Fatalf("Can not create target dir and your home dir(%s), please check your TestChamberBaseDir config in conf/judger.ini was set correctly.", homeDir)
		return
	}

}