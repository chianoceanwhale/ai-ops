package database

import (
	"ai-ops/configs/db"

	log "github.com/sirupsen/logrus"
)

func Setup() {
	dbType := db.DatabaseConfig.Dbtype
	if dbType == "mysql" {
		var mysql = new(Mysql)
		mysql.Setup()
	}
}

type Mysql struct {
}

func (e *Mysql) Setup() {
	var err error

	MysqlConn = e.GetConnect()
	GormDB, err = e.Open(MysqlConn)

	if err != nil {
		log.Panicf("%s connect error %v", DbType, err)
	} else {
		log.Printf("%s connect success!", DbType)
	}

	if GormDB.Error != nil {
		log.Panicf("database error %v", GormDB.Error)
	}

}
