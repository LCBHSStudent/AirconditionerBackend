package Processor

import (
	"encoding/json"
	"fmt"
	"github.com/wxmsummer/airConditioner/server/message"
	"github.com/wxmsummer/airConditioner/server/model"
	"github.com/wxmsummer/airConditioner/server/utils"
	"net"
)

type RoomStateProcessor struct {
	Conn    net.Conn
	RoomNum int
}

// 添加房间状态
func (this *RoomStateProcessor) ProcessRoomStateAdd(msg *message.Message) (err error) {
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
		Cost:        roomStateAdd.Cost,
		Temperature: roomStateAdd.Temperature,
	}
	_, err = model.AddRoomState(roomState)
	if err == nil {
		fmt.Println("add roomState success...")
		roomStateAddRes.Code = 200
	} else {
		fmt.Println("db add roomState err = ", err)
		roomStateAddRes.Code = 500
		roomStateAddRes.Error = err.Error()
	}
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
	tf := &utils.Transfer{Conn: this.Conn}
	err = tf.WritePkg(data)
	return
}

// 查询房间状态
func (this *RoomStateProcessor) ProcessRoomStateQuery(msg *message.Message) (err error) {
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
	roomStateList, err := model.QueryRoomStates(roomNum, startTime, endTime)
	if err == nil {
		roomStateQueryRes.Code = 200
		roomStateQueryRes.RoomStateList = roomStateList
	} else {
		roomStateQueryRes.Code = 500
		roomStateQueryRes.Error = err.Error()
	}
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
	tf := &utils.Transfer{Conn: this.Conn}
	err = tf.WritePkg(data)
	return
}

// 删除房间状态
func (this *RoomStateProcessor) ProcessRoomStateDelete(msg *message.Message) (err error) {
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
	_, err = model.DelRoomStatesByRoomNum(roomNum, startTime, endTime)
	if err == nil {
		roomStateDeleteRes.Code = 200
	} else {
		roomStateDeleteRes.Code = 500
		roomStateDeleteRes.Error = err.Error()
	}
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
	tf := &utils.Transfer{Conn: this.Conn}
	err = tf.WritePkg(data)
	return
}
