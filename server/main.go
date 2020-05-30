package main

import (
	"fmt"
	"github.com/wxmsummer/airConditioner/server/Processor"
	"github.com/wxmsummer/airConditioner/server/db"
	"net"
	"time"
)

// 主控程序，处理和客户端的通讯
func mainProcess(conn net.Conn) {
	defer conn.Close()
	processor := &Processor.MainProcessor{Conn: conn}
	err := processor.Process()
	if err != nil {
		fmt.Println("通讯协程错误，err=", err)
		return
	}
}

// 主函数，初始化连接池、与客户端连接、启动协程
func main() {
	dbName := "mysql"
	dsn := "root:wxm19990516@tcp(127.0.0.1:3306)/airConditioner?charset=utf8"
	maxOpen := 200
	maxIdle := 100
	maxLifeTime := time.Second * 1000
	db.InitMysql(dbName, dsn, maxOpen, maxIdle, maxLifeTime)
	fmt.Println("服务器在8888端口监听...")
	listen, err := net.Listen("tcp", "0.0.0.0:8888")
	if err != nil {
		fmt.Println("net.listen err=", err)
		return
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err=", err)
		}
		fmt.Println("和客户端连接成功...")
		go mainProcess(conn)
	}
}
