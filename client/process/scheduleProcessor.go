package process

import (
	"encoding/json"
	"fmt"
	"github.com/wxmsummer/AirConditioner/common/message"
	"github.com/wxmsummer/AirConditioner/common/utils"
	"net"
)

type ScheduleProcessor struct {
	Conn net.Conn
}

func (sp *ScheduleProcessor) GetServingQueue() (err error) {

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
	msg.Type = message.TypeGetServingQueue

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
	var resMsg message.GetServingQueueRes
	err = json.Unmarshal([]byte(msg.Data), &resMsg)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	fmt.Println(resMsg)

	return
}