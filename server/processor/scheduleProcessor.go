package processor

import (
	"encoding/json"
	"fmt"
	"github.com/wxmsummer/AirConditioner/common/message"
	"github.com/wxmsummer/AirConditioner/common/utils"
	"github.com/wxmsummer/AirConditioner/server/scheduler"
	"net"
)

type ScheduleProcessor struct {
	Conn net.Conn
}

// 查询服务队列
func (sp *ScheduleProcessor) GetServingQueue(msg *message.Message) (err error) {

	var resMsg message.Message
	var getServingQueueRes message.GetServingQueueRes

	servingQueue := scheduler.GetServingQueue()

	getServingQueueRes.Code = 200
	getServingQueueRes.Msg = "查询服务队列成功！"
	getServingQueueRes.ServingQueue = servingQueue

	data, err := json.Marshal(getServingQueueRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeGetServingQueueRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: sp.Conn}
	err = tf.WritePkg(data)
	return
}