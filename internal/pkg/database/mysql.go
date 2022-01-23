package database

import (
	"ai-ops/configs/db"
	"bytes"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DbType   string
	Host     string
	Port     int
	Name     string
	Username string
	Password string
)

var GormDB *gorm.DB
var MysqlConn string

type Database interface {
	Open(conn string) (db *gorm.DB, err error)
	GetConnect() string
}

func (e *Mysql) Open(conn string) (db *gorm.DB, err error) {
	return gorm.Open(mysql.Open(conn), &gorm.Config{})
}

func (e *Mysql) GetConnect() string {

	DbType = db.DatabaseConfig.Dbtype
	Host = db.DatabaseConfig.Host
	Port = db.DatabaseConfig.Port
	Name = db.DatabaseConfig.Name
	Username = db.DatabaseConfig.Username
	Password = db.DatabaseConfig.Password

	var conn bytes.Buffer
	conn.WriteString(Username)
	conn.WriteString(":")
	conn.WriteString(Password)
	conn.WriteString("@tcp(")
	conn.WriteString(Host)
	conn.WriteString(":")
	conn.WriteString(strconv.Itoa(Port))
	conn.WriteString(")")
	conn.WriteString("/")
	conn.WriteString(Name)
	conn.WriteString("?charset=utf8&parseTime=True&loc=Local&timeout=1000ms")
	return conn.String()
}
