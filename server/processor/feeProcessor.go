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

func (fp *FeeProcessor) Add(msg *message.Message) (err error) {
	var feeAdd message.FeeAdd
	err = json.Unmarshal([]byte(msg.Data), &feeAdd)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}
	var resMsg message.Message
	var feeAddRes message.RoomStateAddRes
	fee := &model.Fee{
		RoomNum:   feeAdd.RoomNum,
		StartTime: feeAdd.StartTime,
		EndTime:   feeAdd.EndTime,
		Cost:      feeAdd.Cost,
	}
	err = fp.Orm.AddFee(fee)
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

	resMsg.Type = message.TypeFeeAddRes
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

func (fp *FeeProcessor) Query(msg *message.Message) (err error) {
	var feeQuery message.FeeQuery
	err = json.Unmarshal([]byte(msg.Data), &feeQuery)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}

	var resMsg message.Message
	var feeQueryRes message.FeeQueryRes

	roomNum := feeQuery.RoomNum
	startTime := feeQuery.StartTime
	endTime := feeQuery.EndTime

	fees, err := fp.Orm.QueryFees(roomNum, startTime, endTime)
	if err != nil {
		feeQueryRes.Code = 500
		feeQueryRes.Msg = "fp.Orm.QueryFees err"
	}

	feeQueryRes.Code = 200
	feeQueryRes.Fees = fees

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
	var feeDeleteRes message.FeeDeleteRes

	roomNum := feeDelete.RoomNum
	startTime := feeDelete.StartTime
	endTime := feeDelete.EndTime

	_, err = fp.Orm.DelFees(roomNum, startTime, endTime)
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

	resMsg.Type = message.TypeFeeDeleteRes
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
