package processor

import (
	"encoding/json"
	"fmt"
	"github.com/wxmsummer/AirConditioner/common/message"
	"github.com/wxmsummer/AirConditioner/common/utils"
	"github.com/wxmsummer/AirConditioner/server/model"
	"github.com/wxmsummer/AirConditioner/server/repository"
	"net"
)

type FeeProcessor struct {
	Conn net.Conn
	Orm  *repository.FeeOrm
}

// 插入一条账单信息
func (fp *FeeProcessor) Create(msg *message.Message) (err error) {
	var feeAdd message.FeeAdd
	err = json.Unmarshal([]byte(msg.Data), &feeAdd)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}
	var resMsg message.Message
	var feeAddRes message.NormalRes
	fee := &model.Fee{
		RoomNum:   feeAdd.RoomNum,
		Cost:      feeAdd.Cost,
	}
	err = fp.Orm.Create(fee)
	if err != nil {
		fmt.Println("db add fee err = ", err)
		feeAddRes.Code = 500
		feeAddRes.Msg = "db add fee err = "
	}

	feeAddRes.Code = 200
	feeAddRes.Msg = "add fee success!"

	data, err := json.Marshal(feeAddRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeNormalRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: fp.Conn}
	err = tf.WritePkg(data)
	return
}

// QueryByRoom 根据房间号查询账单
func (fp *FeeProcessor) QueryByRoom(msg *message.Message) (err error) {
	var feeQuery message.FeeQuery
	err = json.Unmarshal([]byte(msg.Data), &feeQuery)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}

	var resMsg message.Message
	var feeQueryRes message.FeeQueryRes

	roomNum := feeQuery.RoomNum

	fee, err := fp.Orm.QueryByRoomNum(roomNum)
	if err != nil {
		feeQueryRes.Code = 500
		feeQueryRes.Msg = "fp.Orm.QueryFees err"
		return
	}

	feeQueryRes.Code = 200
	feeQueryRes.Msg = "查询账单成功！"
	feeQueryRes.Fee = fee

	data, err := json.Marshal(feeQueryRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeFeeQueryRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: fp.Conn}
	err = tf.WritePkg(data)
	return
}

func (fp *FeeProcessor) Delete(msg *message.Message) (err error) {
	var feeDelete message.FeeDelete
	err = json.Unmarshal([]byte(msg.Data), &feeDelete)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}

	var resMsg message.Message
	var feeDeleteRes message.NormalRes

	roomNum := feeDelete.RoomNum


	err = fp.Orm.Delete(roomNum)
	if err != nil {
		feeDeleteRes.Code = 500
		feeDeleteRes.Msg = "fp.Orm.DelFees err"
	}

	feeDeleteRes.Code = 200
	feeDeleteRes.Msg = "del fees success!"

	data, err := json.Marshal(feeDeleteRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeNormalRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: fp.Conn}
	err = tf.WritePkg(data)
	return
}
