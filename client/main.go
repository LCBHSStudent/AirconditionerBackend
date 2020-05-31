package main

import (
	"github.com/wxmsummer/AirConditioner/client/process"
)

func main() {

	// 测试AirProcessor.Create()
	up := &process.AirProcessor{}
	_ = up.Create(302, 1001)

}
