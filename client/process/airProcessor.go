package process

import (

	"encoding/json"
	"fmt"
	"github.com/wxmsummer/AirConditioner/common/message"
	"github.com/wxmsummer/AirConditioner/common/utils"
	"net"
)

type AirProcessor struct {
	Conn net.Conn
}

func (up *AirProcessor) Create(id, roomNum int) (err error) {

	// 1，连接服务器
	conn, err := net.Dial("tcp", "0.0.0.0:8888")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	// 延时关闭
	defer conn.Close()

	// 1，创建一个msg
	var msg message.Message
	msg.Type = message.TypeAirConditionerCreate

	// 2，创建一个smsMsg实例
	var airCreateMsg message.AirConditionerCreate
	airCreateMsg.AirConditioner.Id = id
	airCreateMsg.AirConditioner.RoomNum = roomNum

	// 4，将smsMsg序列化
	data, err := json.Marshal(airCreateMsg)
	if err != nil {
		fmt.Println("airCreateMsg json.Marshal err=", err)
		return
	}

	// 5，赋值data
	msg.Data = string(data)

	// 6，将msg序列化
	data, err = json.Marshal(msg)
	if err != nil {
		fmt.Println("msg json.Marshal err=", err)
		return
	}

	// 实例化一个Transfer
	tf := &utils.Transfer{Conn: conn}

	// 发送消息给服务端
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("sendMsg err=", err)
	}

	// 读取服务端返回的消息
	msg, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg err=", err)
		return
	}

	// 反序列化Data
	var airCreateResMsg message.AirConditionerCreateRes
	err = json.Unmarshal([]byte(msg.Data), &airCreateResMsg)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	fmt.Println(airCreateResMsg)

	return
}
