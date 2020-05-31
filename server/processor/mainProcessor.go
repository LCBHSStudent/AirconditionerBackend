package processor

import (
	"fmt"
	"github.com/wxmsummer/AirConditioner/server/message"
	"github.com/wxmsummer/AirConditioner/server/utils"
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
		err = up.Register(msg)
	case message.TypeUserLogin:
		up := &UserProcessor{Conn: conn}
		err = up.Login(msg)
	case message.TypeUserFindById:
		up := &UserProcessor{Conn: conn}
		err = up.FindById(msg)
	case message.TypeUserFindAll:
		up := &UserProcessor{Conn: conn}
		err = up.FindAll(msg)
	case message.TypeUserUpdate:
		up := &UserProcessor{Conn: conn}
		err = up.Update(msg)

	case message.TypeAirConditionerFindById:
		ap := &AirProcessor{Conn: conn}
		err = ap.FindById(msg)
	case message.TypeAirConditionerFindByRoom:
		ap := &AirProcessor{Conn: conn}
		err = ap.FindByRoom(msg)
	case message.TypeAirConditionerFindAll:
		ap := &AirProcessor{Conn: conn}
		err = ap.FindAll(msg)
	case message.TypeAirConditionerCreate:
		ap := &AirProcessor{Conn: conn}
		err = ap.Create(msg)
	case message.TypeAirConditionerUpdate:
		ap := &AirProcessor{Conn: conn}
		err = ap.Update(msg)

	case message.TypeRoomStateAdd:
		rp := RoomStateProcessor{Conn: conn}
		err = rp.Add(msg)
	case message.TypeRoomStateQuery:
		rp := RoomStateProcessor{Conn: conn}
		err = rp.Query(msg)
	case message.TypeRoomStateDelete:
		rp := RoomStateProcessor{Conn: conn}
		err = rp.Delete(msg)

	case message.TypeFeeAdd:
		fp := FeeProcessor{Conn: conn}
		err = fp.Add(msg)
	case message.TypeFeeQuery:
		fp := FeeProcessor{Conn: conn}
		err = fp.Query(msg)
	case message.TypeFeeDelete:
		fp := FeeProcessor{Conn: conn}
		err = fp.Delete(msg)

	default:
		fmt.Println("消息类型不存在，无法处理...")
	}
	return
}
