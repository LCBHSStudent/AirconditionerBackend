package model

// Report 报表结构体
type Report struct {
	RoomNum     int     `json:"room_num"`      //房间号
	TotalFee    float64 `json:"total_fee"`     //总费用
	TotalPower  float64 `json:"total_power"`   //总耗电量
	CloseNum    int     `json:"close_num"`     //空调的开关次数
	SetParamNum int     `json:"set_param_num"` //空调调整次数
	UsedTime    int     `json:"used_time"`     //使用空调的时长
}
