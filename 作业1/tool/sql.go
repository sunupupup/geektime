package tool

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DBNotExit = errors.New("DB不存在")

var DB *sql.DB
var err error

func init() {
	var err error
	DB, err = sql.Open("mysql", "root:123@tcp(127.0.0.1:3306)/sql_test")
	if err != nil {
		return
	}

	err = DB.Ping()
	if err != nil {
		fmt.Println(err.Error())
	}
	log.Println("DB初始化成功")
}

func GetEngine() *sql.DB {
	return DB
}

func EngineIsExit() (bool, error) {
	if DB == nil {
		return false, DBNotExit
	} else {
		return true, nil
	}
}
