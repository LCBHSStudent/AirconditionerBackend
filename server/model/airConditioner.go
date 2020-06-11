package model

const (
	// 电源开关
	PowerOff = 0
	PowerOn  = 1

	// 空调模式
	ModeCold      = 0 // 制冷
	ModeHot       = 1 // 制热
	ModeWind      = 2 // 送风
	ModeDry       = 3 // 干燥
	ModeSleep     = 4 // 睡眠
	ModeSwingFlap = 5
	ModeBreath    = 6

	// 风速
	WindAuto = 0
	WindLow  = 1
	WindMid  = 2
	WindHigh = 3

	// 空调初始化温度
	InitTemperature = 25
)

// 空调数据结构
type AirConditioner struct {
	RoomNum         int     `json:"room_num"`         // 空调所在房间号，默认一个房间一个空调(中央空调形式)
	Power           int     `json:"power"`            // 电源开关：0关 1开
	Mode            int     `json:"mode"`             // 模式
	WindLevel       int     `json:"wind_level"`       // 风速
	Temperature     float64 `json:"temperature"`      // 温度
	RoomTemperature float64 `json:"room_temperature"` // 室温
	TotalPower      float64 `json:"total_power"`      // 该次入住的总耗电量
	StartWind       []int64 `json:"start_wind"`       // 开始送风时间，时间戳格式
	StopWind        []int64 `json:"stop_wind"`        // 停止送风时间
	OpenTime        []int64 `json:"open_time"`        // 开机时间，数组
	CloseTime       []int64 `json:"close_time"`       // 关机时间，数组
	SetParamNum     int     `json:"set_param_num"`    // 调整次数，用于报表展示
}
