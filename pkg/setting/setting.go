package setting

import (
	"github.com/go-ini/ini"
	"log"
	"os"
	"path/filepath"
)

const FRONTENDOJUDGER = "FrontEndOJudger"

type Judger struct {
	JudgerSum          int
	BufferSum		int
	SleepTime          int

	TestChamberSwitch bool
	TestChamberBaseDir string
	TestChamberAddr string
	TestChamberPort string
	HttpJudgerPort string
}

var JudgerSetting = &Judger{}

type Frontend struct {
	EnableWebsocket bool
	WebsocketScheme string
	WebsocketHost string
	WebsocketPath string
	WebsocketToken string
}
var FrontEndSetting = &Frontend{}

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

type File struct {
	FileToolType string
	FileBaseDir string
	CloudFileBaseUrl string
}

var FileSetting = &File{}

var cfg *ini.File

func Setup() {
	var err error
	cfg, err = ini.Load("conf/judger.ini")
	if err != nil {
		log.Fatalf("Judger setup failed, place check [conf/judger.ini] file exist error:%v", err)
	}

	mapTo("database", DatabaseSetting)
	mapTo("judger", JudgerSetting)
	mapTo("frontend", FrontEndSetting)
	mapTo("file", FileSetting)
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
	//checkLogger()
}


func checkAndFixDirExists(targetDir string, suffix string) string {
	_, err := os.Stat(targetDir)
	if err == nil || os.IsExist(err) {
		return targetDir
	}
	// try create dir if not exist
	err = os.MkdirAll(targetDir, 0777)
	if err == nil {
		return targetDir
	}

	// failed to create dir but try user home dir
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Can not create target dir and your home dir, please check your config in conf/*.ini was set correctly.")
	}
	newDir := filepath.Join(homeDir, suffix)
	// try create dir if not exist
	err = os.MkdirAll(newDir, 0777)
	if err == nil {
		return newDir
	}
	log.Fatalf("Can not create home dir, please check your config in conf/*.ini was set correctly.")
	return ""
}


func checkJudger() {
	// check judger TestChamberBaseDir
	JudgerSetting.TestChamberBaseDir = checkAndFixDirExists(JudgerSetting.TestChamberBaseDir, "test_submit")
	FileSetting.FileBaseDir = checkAndFixDirExists(FileSetting.FileBaseDir, "static/file")
}