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

func (up *AirProcessor) Create(roomNum int) (err error) {

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
	var resMsg message.NormalRes
	err = json.Unmarshal([]byte(msg.Data), &resMsg)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	fmt.Println(resMsg)

	return
}

func (up *AirProcessor) PowerOn(powerOn message.AirConditionerOn) (err error) {

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
	msg.Type = message.TypeAirConditionerOn

	// 4，将smsMsg序列化
	data, err := json.Marshal(powerOn)
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
	var resMsg message.NormalRes
	err = json.Unmarshal([]byte(msg.Data), &resMsg)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	fmt.Println(resMsg)

	return
}

func (up *AirProcessor) SetParam(setParam message.AirConditionerSetParam) (err error) {

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
	msg.Type = message.TypeAirConditionerSetParam

	// 4，将smsMsg序列化
	data, err := json.Marshal(setParam)
	if err != nil {
		fmt.Println("setParam json.Marshal err=", err)
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
	var resMsg message.NormalRes
	err = json.Unmarshal([]byte(msg.Data), &resMsg)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	fmt.Println(resMsg)

	return
}

func (up *AirProcessor) PowerOff(powerOff message.AirConditionerOff) (err error) {

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
	msg.Type = message.TypeAirConditionerOff

	// 4，将smsMsg序列化
	data, err := json.Marshal(powerOff)
	if err != nil {
		fmt.Println("powerOff json.Marshal err=", err)
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
	var resMsg message.NormalRes
	err = json.Unmarshal([]byte(msg.Data), &resMsg)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	fmt.Println(resMsg)

	return
}

func (up *AirProcessor) WatchAir() (err error) {

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
	msg.Type = message.TypeAirConditionerFindAll

	// 6，将msg序列化
	data, err := json.Marshal(msg)
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
	var resMsg message.AirConditionerFindAllRes
	err = json.Unmarshal([]byte(msg.Data), &resMsg)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	fmt.Println(resMsg)

	return
}