package processor

import (
	"encoding/json"
	"fmt"
	"github.com/wxmsummer/AirConditioner/server/message"
	"github.com/wxmsummer/AirConditioner/server/repository"
	"github.com/wxmsummer/AirConditioner/server/utils"
	"net"
)

type AirProcessor struct {
	Conn net.Conn
	Orm  *repository.AirConditionerOrm
}

func (ap *AirProcessor) FindById(msg *message.Message) (err error) {
	var airConditionerFindById message.AirConditionerFindById
	err = json.Unmarshal([]byte(msg.Data), &airConditionerFindById)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}

	var resMsg message.Message
	var airConditionerFindByIdRes message.AirConditionerFindByIdRes

	id := airConditionerFindById.Id
	airConditioner, err := ap.Orm.FindByID(id)
	if err != nil {
		fmt.Println("ap.Orm.FindByID(id) err=", err)
		return err
	}

	airConditionerFindByIdRes.Code = 200
	airConditionerFindByIdRes.Msg = "查询成功！"
	airConditionerFindByIdRes.AirConditioner = airConditioner

	data, err := json.Marshal(airConditionerFindByIdRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeAirConditionerFindByIdRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: ap.Conn}
	err = tf.WritePkg(data)
	return
}

// 查询空调状态，按房间号查询
func (ap *AirProcessor) FindByRoom(msg *message.Message) (err error) {
	var airConditionerFindByRoom message.AirConditionerFindByRoom
	err = json.Unmarshal([]byte(msg.Data), &airConditionerFindByRoom)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}

	var resMsg message.Message
	var airConditionerFindByRoomRes message.AirConditionerFindByRoomRes

	roomNum := airConditionerFindByRoom.RoomNum
	airConditioners, err := ap.Orm.FindByRoom(roomNum)
	if err != nil {
		fmt.Println("ap.Orm.FindByRoom err=", err)
		return err
	}

	airConditionerFindByRoomRes.Code = 200
	airConditionerFindByRoomRes.Msg = "查询成功！"
	airConditionerFindByRoomRes.AirConditioners = airConditioners

	data, err := json.Marshal(airConditionerFindByRoomRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeAirConditionerFindByRoomRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: ap.Conn}
	err = tf.WritePkg(data)
	return
}

func (ap *AirProcessor) FindAll(msg *message.Message) (err error) {
	var airConditionerFindAll message.AirConditionerFindAll
	err = json.Unmarshal([]byte(msg.Data), &airConditionerFindAll)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}

	var resMsg message.Message
	var airConditionerFindAllRes message.AirConditionerFindAllRes

	airConditioners, err := ap.Orm.FindAll()
	if err != nil {
		fmt.Println("ap.Orm.FindAll() err=", err)
		return err
	}

	airConditionerFindAllRes.Code = 200
	airConditionerFindAllRes.Msg = "查询成功！"
	airConditionerFindAllRes.AirConditioners = airConditioners

	data, err := json.Marshal(airConditionerFindAllRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeAirConditionerFindAllRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: ap.Conn}
	err = tf.WritePkg(data)
	return
}

func (ap *AirProcessor) Create(msg *message.Message) (err error) {
	var airConditionerCreate message.AirConditionerCreate
	err = json.Unmarshal([]byte(msg.Data), &airConditionerCreate)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}

	var resMsg message.Message
	var airConditionerCreateRes message.AirConditionerCreateRes

	airConditioner := airConditionerCreate.AirConditioner
	err = ap.Orm.Create(&airConditioner)
	if err != nil {
		fmt.Println("ap.Orm.Create(&airConditioner) err=", err)
		return
	}

	airConditionerCreateRes.Code = 200
	airConditionerCreateRes.Msg = "创建成功！"

	data, err := json.Marshal(airConditionerCreateRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeAirConditionerCreateRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: ap.Conn}
	err = tf.WritePkg(data)
	return
}

func (ap *AirProcessor) Update(msg *message.Message) (err error) {
	var airConditionerUpdate message.AirConditionerUpdate
	err = json.Unmarshal([]byte(msg.Data), &airConditionerUpdate)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}

	var resMsg message.Message
	var airConditionerUpdateRes message.AirConditionerUpdateRes

	airConditioner := airConditionerUpdate.AirConditioner
	err = ap.Orm.Update(airConditioner)
	if err != nil {
		fmt.Println("ap.Orm.Update(airConditioner) err=", err)
		return
	}

	airConditionerUpdateRes.Code = 200
	airConditionerUpdateRes.Msg = "更新成功！"

	data, err := json.Marshal(airConditionerUpdateRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeAirConditionerUpdateRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: ap.Conn}
	err = tf.WritePkg(data)
	return
}
