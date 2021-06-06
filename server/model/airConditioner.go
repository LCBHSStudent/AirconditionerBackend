package model

const (
	// 电源开关
	PowerOn  = "on"
	PowerOff = "off"
	
	// 空调模式
	ModeCold      = "cold" // 制冷
	ModeHot       = "hot" // 制热
	ModeWind      = "wind" // 送风
	ModeDry       = "dry" // 干燥
	ModeSleep     = "sleep" // 睡眠
	ModeSwingFlap = "swing"
	ModeBreath    = "breath"

	// 风速
	WindAuto = "auto"
	WindLow  = "low"
	WindMid  = "mid"
	WindHigh = "high"

	// 空调初始化温度
	InitTemperature = 25
)

var LevelMap = map[string]int{
	"stop":0,
	"low":1,
	"mid":2,
	"high":3,
}

// 空调数据结构
type AirConditioner struct {
	RoomNum         int     `json:"room_num" gorm:"not null;unique"` // 空调所在房间号，默认一个房间一个空调(中央空调形式)
	Power           string  `json:"power"`                           // 电源开关：on开 off关
	Mode            string  `json:"mode"`                            // 模式
	WindLevel       string  `json:"wind_level"`                      // 风速
	Temperature     float64 `json:"temperature"`                     // 温度
	RoomTemperature float64 `json:"room_temperature"`                // 室温
	TotalPower      float64 `json:"total_power"`                     // 该次入住的总耗电量
	TotalFee		float64	`json:"total_fee"`						 // 该次入住的总费用
	StartWind       string  `json:"start_wind" gorm:"type:TEXT"`     // 开始送风时间，时间戳格式
	StopWind        string  `json:"stop_wind" gorm:"type:TEXT"`      // 停止送风时间
	OpenTime        string  `json:"open_time" gorm:"type:TEXT"`      // 开机时间，数组，如 [1,2,3,4] ，由于mysql不支持切片类型，转换为string存储
	CloseTime       string  `json:"close_time" gorm:"type:TEXT"`     // 关机时间，数组
	SetParamNum     int     `json:"set_param_num"`                   // 调整次数，用于报表展示
}

