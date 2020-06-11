package processor

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/wxmsummer/AirConditioner/common/message"
	"github.com/wxmsummer/AirConditioner/common/utils"
	"github.com/wxmsummer/AirConditioner/server/repository"
	"io"
	"net"
)

type MainProcessor struct {
	Conn net.Conn
	Db   *gorm.DB
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
	userOrm := &repository.UserOrm{Db: this.Db}
	airOrm := &repository.AirConditionerOrm{Db: this.Db}
	feeOrm := &repository.FeeOrm{Db: this.Db}

	switch msg.Type {
	case message.TypeUserRegister:
		up := &UserProcessor{Conn: conn, Orm: userOrm}
		err = up.Register(msg)
	case message.TypeUserLogin:
		up := &UserProcessor{Conn: conn, Orm: userOrm}
		err = up.Login(msg)
	case message.TypeUserFindById:
		up := &UserProcessor{Conn: conn, Orm: userOrm}
		err = up.FindById(msg)
	case message.TypeUserFindAll:
		up := &UserProcessor{Conn: conn, Orm: userOrm}
		err = up.FindAll(msg)
	case message.TypeUserUpdate:
		up := &UserProcessor{Conn: conn, Orm: userOrm}
		err = up.Update(msg)

	case message.TypeAirConditionerFindByRoom:
		ap := &AirProcessor{Conn: conn, Orm: airOrm}
		err = ap.FindByRoom(msg)
	case message.TypeAirConditionerFindAll:
		ap := &AirProcessor{Conn: conn, Orm: airOrm}
		err = ap.FindAll(msg)
	case message.TypeAirConditionerCreate:
		ap := &AirProcessor{Conn: conn, Orm: airOrm}
		err = ap.Create(msg)
	case message.TypeAirConditionerUpdate:
		ap := &AirProcessor{Conn: conn, Orm: airOrm}
		err = ap.Update(msg)

	case message.TypeAirConditionerOn:
		ap := &AirProcessor{Conn: conn, Orm: airOrm}
		err = ap.PowerOn(msg)
	case message.TypeAirConditionerOff:
		ap := &AirProcessor{Conn: conn, Orm: airOrm}
		err = ap.PowerOff(msg)
	case message.TypeAirConditionerSetParam:
		ap := &AirProcessor{Conn: conn, Orm: airOrm}
		err = ap.SetParam(msg)
	case message.TypeAirConditionerStopWind:
		ap := &AirProcessor{Conn: conn, Orm: airOrm}
		err = ap.StopWind(msg)
	case message.TypeGetReport:
		ap := &AirProcessor{Conn: conn, Orm: airOrm}
		err = ap.GetReport(msg)

	case message.TypeFeeAdd:
		fp := FeeProcessor{Conn: conn, Orm: feeOrm}
		err = fp.Add(msg)
	case message.TypeFeeQuery:
		fp := FeeProcessor{Conn: conn, Orm: feeOrm}
		err = fp.QueryByRoom(msg)
	case message.TypeFeeDelete:
		fp := FeeProcessor{Conn: conn, Orm: feeOrm}
		err = fp.Delete(msg)

	default:
		fmt.Println("消息类型不存在，无法处理...")
	}
	return
}
