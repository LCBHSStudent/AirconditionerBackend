package main

import (
	"fmt"
	"github.com/wxmsummer/AirConditioner/server/database"
	"github.com/wxmsummer/AirConditioner/server/processor"
	"log"
	"net"
)

// 主控程序，处理和客户端的通讯
func mainProcess(conn net.Conn) {
	defer conn.Close()
	mainProcessor := &processor.MainProcessor{Conn: conn}
	err := mainProcessor.Process()
	if err != nil {
		fmt.Println("通讯协程错误，err=", err)
		return
	}
}

// 主函数，初始化连接池、与客户端连接、启动协程
func main() {

	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("connection error : %v \n", err)
	}
	defer db.Close()

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
