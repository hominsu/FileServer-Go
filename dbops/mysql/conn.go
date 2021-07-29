package mysql

import "database/sql"

var (
	dbConn *sql.DB
	err    error
)

// 初始化，只执行一次，用一个全局变量保存，确保性能
func init() {
	dbConn, err = sql.Open("mysql", "root:root@tcp(localhost:3310)/file_server?charset=utf8")
	if err != nil {
		panic(err.Error())
	}

	dbConn.SetMaxIdleConns(1000)
	err = dbConn.Ping()
	if err != nil {
		panic(err.Error())
	}
}

func DBConn() *sql.DB {
	return dbConn
}
