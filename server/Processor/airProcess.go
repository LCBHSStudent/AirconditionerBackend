package Processor

import (
	"encoding/json"
	"fmt"
	"github.com/wxmsummer/airConditioner/server/message"
	"github.com/wxmsummer/airConditioner/server/model"
	"github.com/wxmsummer/airConditioner/server/utils"
	"net"
)

type AirProcessor struct {
	Conn    net.Conn
	RoomNum int
}

// 修改空调状态
func (this *AirProcessor) ProcessAirModify(msg *message.Message) (err error) {
	var airModify message.AirModify
	err = json.Unmarshal([]byte(msg.Data), &airModify)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}
	var resMsg message.Message
	var airModifyRes message.AirModifyRes
	airConditioner := &model.AirConditioner{
		Number:      airModify.Number,
		Power:       airModify.Power,
		Mode:        airModify.Mode,
		WindLevel:   airModify.WindLevel,
		Temperature: airModify.Temperature,
	}
	_, err = model.UpdateAirConditioner(airConditioner) // 到数据库去更新空调状态
	if err == nil {
		airModifyRes.Code = 200
	} else {
		airModifyRes.Code = 500
		airModifyRes.Error = err.Error()
	}
	data, err := json.Marshal(airModifyRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}
	resMsg.Type = message.TypeAirModifyRes
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

// 查询空调状态，如果传入的房间号为0，则查询所有空调信息，否则按房间号查询
func (this *AirProcessor) ProcessAirQuery(msg *message.Message) (err error) {
	var airQuery message.AirQuery
	err = json.Unmarshal([]byte(msg.Data), &airQuery)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}
	var resMsg message.Message
	var airQueryRes message.AirQueryRes
	roomNum := airQuery.RoomNumber
	if roomNum != 0 { // 查询单个空调的信息
		air := model.QueryAirWithNumber(roomNum)
		var airList []model.AirConditioner
		airList = append(airList, air)
		airQueryRes.Code = 200
		airQueryRes.AirList = airList
	} else if roomNum == 0 { // 查询所有空调的信息
		airList, err := model.QueryAllAirConditioners()
		if err == nil {
			airQueryRes.Code = 200
			airQueryRes.AirList = airList
		} else {
			airQueryRes.Code = 500
			airQueryRes.Error = err.Error()
		}
	}
	data, err := json.Marshal(airQueryRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}
	resMsg.Type = message.TypeAirQueryRes
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
