package model

// 费用结构体，记录房间号和费用
type Fee struct {
	RoomNum int     `json:"room_num"` // 房间号
	Cost    float64 `json:"cost"`     // 此次入住的费用
}
