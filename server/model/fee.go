package model

// 费用结构体，记录房间号、开始时间、结束时间，以及这段时间内的费用
type Fee struct {
	RoomNum   int     `json:"room_num"`   // 房间号
	StartTime int64   `json:"start_time"` // 开始时间
	EndTime   int64   `json:"end_time"`   // 结束时间
	Cost      float64 `json:"cost"`       // 费用
}
