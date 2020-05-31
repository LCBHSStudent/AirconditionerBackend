package model

const (
	// 电源开关
	PowerOff = 0
	PowerOn  = 1

	// 空调模式
	ModeCold      = 0
	ModeHot       = 1
	ModeWind      = 2
	ModeDry       = 3
	ModeSleep     = 4
	ModeSwingFlap = 5
	ModeBreath    = 6

	// 风速
	WindAuto = 0
	WindLow  = 1
	WindMid  = 2
	WindHigh = 3
)

// 空调数据结构
type AirConditioner struct {
	Id          int     `json:"id"`          // 空调编号
	RoomNum     int     `json:"room_num"`    // 空调所在房间号，一个房间可能有多个空调
	Power       int     `json:"power"`       // 电源开关：0关 1开
	Mode        int     `json:"mode"`        // 模式
	WindLevel   int     `json:"wind_level"`  // 风速
	Temperature float64 `json:"temperature"` // 温度
}
