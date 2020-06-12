package model

// 空调数据结构
type AirConditioner struct {
	RoomNum         int     `json:"room_num"`         // 空调所在房间号，默认一个房间一个空调(中央空调形式)
	Power           int     `json:"power"`            // 电源开关：0关 1开
	Mode            int     `json:"mode"`             // 模式
	WindLevel       int     `json:"wind_level"`       // 风速
	Temperature     float64 `json:"temperature"`      // 温度
	RoomTemperature float64 `json:"room_temperature"` // 室温
	TotalPower      float64 `json:"total_power"`      // 该次入住的总耗电量
	StartWind       string  `json:"start_wind"`       // 开始送风时间，时间戳格式
	StopWind        string  `json:"stop_wind"`        // 停止送风时间
	OpenTime        string  `json:"open_time"`        // 开机时间，数组，如 [1,2,3,4] ，由于mysql不支持切片类型，转换为string存储
	CloseTime       string  `json:"close_time"`       // 关机时间，数组
	SetParamNum     int     `json:"set_param_num"`    // 调整次数，用于报表展示
}
