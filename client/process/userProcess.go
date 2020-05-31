package process

import (
	"fmt"
	"net"
)

type UserProcess struct {
	// 暂时不需要任何字段
}

// 登录函数
func (this *UserProcess) Login(userID int, userPwd string) (err error) {

	// 定协议。。
	// 1，连接服务器
	conn, err := net.Dial("tcp", "0.0.0.0:8888")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	// 延时关闭
	defer conn.Close()

	// 需要起一个协程，用于保持和服务器端的通讯，如果服务器有数据推送给客户端，
	// 则接收并显示在客户端终端
	go serverProcessMsg(conn)

	return
}
