package database

import (
	"errors"
	"fmt"
	"gin-test/util/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	gormDB *gorm.DB
)

const (
	UserTableName    = "user"
	SafeJobTableName = "safe_job"
	SafeLogTableName = "safe_log"

	Undeleted = 0
	Deleted   = 1
)

func Init() error {
	if gormDB != nil {
		log.Log(fmt.Sprintf("database exists error\n"))
		return errors.New("gormDB not nil")
	}
	host := "127.0.0.1"
	port := "3306"
	database := "bjfu"
	username := "root"
	password := "Baidu@2022"
	charset := "utf8"
	address := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset,
	)
	var err error
	gormDB, err = gorm.Open(mysql.Open(address), &gorm.Config{})
	if err != nil {
		log.Log(fmt.Sprintf("connect database failed: %s\n", err))
		return err
	}
	return nil
}

func GetDB() *gorm.DB {
	// return gormDB.Debug()
	return gormDB
}

type BaseDao struct {
	Engine *gorm.DB
}

func (p *BaseDao) GetDB() *gorm.DB {
	return p.Engine
}

func (p *BaseDao) Transaction(db *gorm.DB) {
	p.Engine = db
}

func IsError(err error) bool {
	if err != nil && err != gorm.ErrRecordNotFound {
		return true
	}

	return false
}

func IsRecordNotFound(err error) bool {
	return err == gorm.ErrRecordNotFound
}
