package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var db *sql.DB

// 初始化mysql、创建连接池
func InitMysql(name, dsn string, maxOpen, maxIdle int, maxLifeTime time.Duration) {
	db1, err := sql.Open(name, dsn)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		db = db1
		db.SetMaxOpenConns(maxOpen)
		db.SetMaxIdleConns(maxIdle)
		db.SetConnMaxLifetime(maxLifeTime)
		CreateTableWithUser()
		CreateTableWithAirConditioner()
		CreateTableWithRoomStates()
	}
}

// 操作数据库
func ModifyDB(sql string, args ...interface{}) (int64, error) {
	result, err := db.Exec(sql, args...)
	if err != nil {
		fmt.Println("db.Exec err =", err)
		return 0, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		fmt.Println("result.RowsAffected err=", err)
		return 0, err
	}
	return count, err
}

// 查询单行数据
func QueryRowDB(sql string) *sql.Row {
	return db.QueryRow(sql)
}

// 查询多行数据
func QueryRowsDB(sql string) (*sql.Rows, error) {
	return db.Query(sql)
}

// 创建用户表
func CreateTableWithUser() {
	sqlStr := `create table if not exists users(
	id int(4) primary key auto_increment not null,
	roomNum int(4),
	privilege int(1),
	password varchar(64),
	checkIn int(10),
	checkOut int(10)
	);`
	_, _ = ModifyDB(sqlStr)
}

// 创建空调状态表
func CreateTableWithAirConditioner() {
	sqlStr := `create table if not exists airs(
	number int(4) primary key auto_increment not null,
	power int(1),
	mode int(1),
	windLevel int(1),
	temperature float(2,2)
	);`
	_, _ = ModifyDB(sqlStr)
}

// 创建房间状态表，存储温度、耗电量和费用
func CreateTableWithRoomStates() {
	sqlStr := `create table if not exists roomStates(
	id int(11) primary key auto_increment not null,
	roomNum int(4),
	startTime int(10),
	endTime int(10),
	power float(5,2),
	cost float(5, 2),
	temper float(2,2)
	);`
	_, _ = ModifyDB(sqlStr)
}
