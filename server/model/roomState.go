package model

// 房间状态结构体，记录房间号、开始时间、结束时间，以及这段时间内的温度和耗电量
type RoomState struct {
	RoomNum     int     `json:"room_num"`    // 房间号
	StartTime   int64   `json:"start_time"`  // 开始时间
	EndTime     int64   `json:"end_time"`    // 结束时间
	Power       float64 `json:"power"`       // 耗电量
	Temperature float64 `json:"temperature"` // 温度
}
