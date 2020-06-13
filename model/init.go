package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"seabase/extend/conf"
	"time"
)

type BaseModel struct {
	CreatedAt time.Time `gorm:"column:created_at;type:datetime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
}

var DB *gorm.DB

func Init() {
	var err error
	fmt.Println("----", conf.DBConf.DBType, conf.DBConf.ConnStr)
	DB, err = gorm.Open(conf.DBConf.DBType, conf.DBConf.ConnStr)
	if err != nil {
		fmt.Println("Mysql server connect error %v", err)
		time.Sleep(10 * time.Second)
		DB, err = gorm.Open(conf.DBConf.DBType, conf.DBConf.ConnStr)
		if err != nil {
			panic(err.Error())
		}
	}
	if DB.Error != nil {
		fmt.Printf("database error %v", DB.Error)
	}
	DB.LogMode(conf.DBConf.Debug)
	DB.SingularTable(true)
	DB.DB().SetConnMaxLifetime(time.Minute * 10)
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)
}
