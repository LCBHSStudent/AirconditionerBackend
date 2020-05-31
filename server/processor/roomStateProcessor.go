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

type RoomStateProcessor struct {
	Conn net.Conn
	Orm  *repository.RoomStateOrm
}

// 添加房间状态
func (rp *RoomStateProcessor) Add(msg *message.Message) (err error) {
	var roomStateAdd message.RoomStateAdd
	err = json.Unmarshal([]byte(msg.Data), &roomStateAdd)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}
	var resMsg message.Message
	var roomStateAddRes message.RoomStateAddRes
	roomState := &model.RoomState{
		RoomNum:     roomStateAdd.RoomNum,
		StartTime:   roomStateAdd.StartTime,
		EndTime:     roomStateAdd.EndTime,
		Power:       roomStateAdd.Power,
		Temperature: roomStateAdd.Temperature,
	}
	err = rp.Orm.AddRoomState(roomState)
	if err != nil {
		fmt.Println("db add roomState err = ", err)
		roomStateAddRes.Code = 500
		roomStateAddRes.Msg = "db add roomState err = "
	}

	roomStateAddRes.Code = 200
	roomStateAddRes.Msg = "add roomState success!"

	data, err := json.Marshal(roomStateAddRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeRoomStateAddRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: rp.Conn}
	err = tf.WritePkg(data)
	return
}

// 查询房间状态
func (rp *RoomStateProcessor) Query(msg *message.Message) (err error) {
	var roomStateQuery message.RoomStateQuery
	err = json.Unmarshal([]byte(msg.Data), &roomStateQuery)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}

	var resMsg message.Message
	var roomStateQueryRes message.RoomStateQueryRes

	roomNum := roomStateQuery.RoomNum
	startTime := roomStateQuery.StartTime
	endTime := roomStateQuery.EndTime

	roomStates, err := rp.Orm.QueryRoomStates(roomNum, startTime, endTime)
	if err != nil {
		roomStateQueryRes.Code = 500
		roomStateQueryRes.Msg = "rp.Orm.QueryRoomStates err"
	}

	roomStateQueryRes.Code = 200
	roomStateQueryRes.RoomStates = roomStates

	data, err := json.Marshal(roomStateQueryRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeRoomStateQueryRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: rp.Conn}
	err = tf.WritePkg(data)
	return
}

// 删除房间状态
func (rp *RoomStateProcessor) Delete(msg *message.Message) (err error) {
	var roomStateDelete message.RoomStateDelete
	err = json.Unmarshal([]byte(msg.Data), &roomStateDelete)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}

	var resMsg message.Message
	var roomStateDeleteRes message.RoomStateDeleteRes

	roomNum := roomStateDelete.RoomNum
	startTime := roomStateDelete.StartTime
	endTime := roomStateDelete.EndTime

	_, err = rp.Orm.DelRoomStates(roomNum, startTime, endTime)
	if err != nil {
		roomStateDeleteRes.Code = 500
		roomStateDeleteRes.Msg = "rp.Orm.DelRoomStates err"
	}

	roomStateDeleteRes.Code = 200
	roomStateDeleteRes.Msg = "del roomStates success!"

	data, err := json.Marshal(roomStateDeleteRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeRoomStateDeleteRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: rp.Conn}
	err = tf.WritePkg(data)
	return
}
