package setting

import (
	"github.com/go-ini/ini"
	"log"
)

type Judger struct {
	TestChamberBaseDir string
	JudgerSum int
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
		log.Fatalf("Judger setup failed, place check [conf/judger.init] file exist error:%v", err)
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
