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

// 默认费率
const DefaultFeeRate = 1

// Detail 详单结构体
// 房间号、开始送风时间、结束送风时间、送风时长、风速、费率、费用。
type Detail struct {
	RoomNum      	int   	 `json:"room_num"`
	StartWindList	[]int64	 `json:"start_wind_list"`
	StoptWindList	[]int64	 `json:"stop_wind_list"`
	TotalWindTime   int64	 `json:"total_wind_time"`
	WindLevel		string 	 `json:"wind_level"`
	TotalPower      float64  `json:"total_power"`
	FeeRate			float64  `json:"fee_rate"`
	TotalFee		float64	 `json:"total_fee"`
}