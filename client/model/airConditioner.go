package model


// 空调数据结构
type AirConditioner struct {
	Id          int     `json:"id"`          // 空调编号
	RoomNum     int     `json:"room_num"`    // 空调所在房间号，一个房间可能有多个空调
	Power       int     `json:"power"`       // 电源开关：0关 1开
	Mode        int     `json:"mode"`        // 模式
	WindLevel   int     `json:"wind_level"`  // 风速
	Temperature float64 `json:"temperature"` // 温度
}
