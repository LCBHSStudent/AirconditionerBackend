package main

import (
	"github.com/wxmsummer/AirConditioner/client/process"
	"github.com/wxmsummer/AirConditioner/common/message"
)

func main() {

	// 测试AirProcessor.Create()
	up := &process.AirProcessor{}
	// _ = up.Create(1002)

	// powerOn := message.AirConditionerOn{
	// 	RoomNum:     1001,
	// 	Mode:        1,
	// 	WindLevel:   1,
	// 	Temperature: 25.5,
	// 	OpenTime:    1591873406,
	// }
	// _ = up.PowerOn(powerOn)

	setParam := message.AirConditionerSetParam {
		RoomNum: 1001,
		Mode: 2,
		WindLevel: 2,
		Temperature: 26,
	}

	_ = up.SetParam(setParam)
}
