package Processor

import (
	"fmt"
	"github.com/wxmsummer/airConditioner/server/message"
	"github.com/wxmsummer/airConditioner/server/utils"
	"io"
	"net"
)

type MainProcessor struct {
	Conn net.Conn
}

// Process用于监听并处理客户端发来的消息
func (this *MainProcessor) Process() (err error) {
	conn := this.Conn
	for {
		tf := &utils.Transfer{Conn: conn}
		msg, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出...服务端也退出..")
				return err
			} else {
				fmt.Println("readPkg err=", err)
				return err
			}
		}
		err = this.serverProcessMsg(&msg)
		if err != nil {
			return err
		}
		fmt.Println("msg=", msg)
	}
}

// 编写一个serverProcessMsg函数
// 功能：根据客户端发送的消息种类不同，决定调用哪个函数来处理
func (this *MainProcessor) serverProcessMsg(msg *message.Message) (err error) {
	conn := this.Conn
	switch msg.Type {
	case message.TypeUserRegister:
		up := &UserProcessor{Conn: conn}
		err = up.ProcessUserRegister(msg)
	case message.TypeUserLogin:
		up := &UserProcessor{Conn: conn}
		err = up.ProcessUserLogin(msg)
	case message.TypeUserModify:
		up := &UserProcessor{Conn: conn}
		err = up.ProcessUserModify(msg)
	case message.TypeUserQuery:
		up := &UserProcessor{Conn: conn}
		err = up.ProcessUserQuery(msg)

	case message.TypeAirModify:
		ap := &AirProcessor{Conn: conn}
		err = ap.ProcessAirModify(msg)
	case message.TypeAirQuery:
		ap := &AirProcessor{Conn: conn}
		err = ap.ProcessAirQuery(msg)

	case message.TypeRoomStateAdd:
		rp := RoomStateProcessor{Conn: conn}
		err = rp.ProcessRoomStateAdd(msg)
	case message.TypeRoomStateQuery:
		rp := RoomStateProcessor{Conn: conn}
		err = rp.ProcessRoomStateQuery(msg)
	case message.TypeRoomStateDelete:
		rp := RoomStateProcessor{Conn: conn}
		err = rp.ProcessRoomStateDelete(msg)
	default:
		fmt.Println("消息类型不存在，无法处理...")
	}
	return
}
