package process

import (
	"fmt"
	"github.com/wxmsummer/AirConditioner/common/message"
	"github.com/wxmsummer/AirConditioner/common/utils"
	"net"
)

// 和服务端保持通讯
func serverProcessMsg(conn net.Conn) {

	// 创建一个Transfer，让其不停地读取服务器发送的消息
	tf := &utils.Transfer{Conn: conn}
	for {
		fmt.Println("客户端正在等待服务器发送的消息...")
		msg, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("tf.ReadPkg err=", err)
			return
		}
		// 如果读取到消息，则进行下一步处理
		switch msg.Type {
		case message.TypeUserFindAllRes:
		default:
			fmt.Println("服务器端返回了未知的消息类型...")
		}
		fmt.Printf("Msg=%v\n", msg)
	}
}
