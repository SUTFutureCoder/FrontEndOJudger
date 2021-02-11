package models

import (
	"FrontEndOJudger/pkg/setting"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type Model struct {
	// ID 自增ID
	ID uint64 `json:"id"`
	// Status 实验室状态
	Status int `json:"status"`
	// CreatorId 创建人Id
	CreatorId uint64 `json:"creator_id"`
	// Creator 创建人
	Creator string `json:"creator"`
	// CreateTime 创建时间
	CreateTime int64 `json:"create_time"`
	// UpdateTime 修改时间
	UpdateTime int `json:"update_time"`
}

const (
	STATUS_CONSTRUCTING = -2
	STATUS_ALL = -1
	STATUS_DISABLE = 0
	STATUS_ENABLE = 1
)

var DB *sql.DB

func Setup() {
	var err error
	DB, err = sql.Open(setting.DatabaseSetting.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name))
	if err != nil {
		log.Fatalf("[FATAL] Database init error[%v]", err)
		return
	}

	DB.SetMaxIdleConns(setting.DatabaseSetting.MaxIdleConns)
	DB.SetMaxOpenConns(setting.DatabaseSetting.MaxOpenConns)
	err = DB.Ping()
	if err != nil {
		log.Fatalf("[FATAL] Database ping error[%v]", err)
	}
}
