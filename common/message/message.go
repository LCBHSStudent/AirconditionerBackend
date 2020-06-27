package message

import "github.com/wxmsummer/AirConditioner/server/model"

// 消息类型
const (
	TypeNormalRes = "NormalRes"

	TypeUserRegister    = "UserRegister"
	TypeUserLogin       = "UserLogin"
	TypeUserFindById    = "UserFindById"
	TypeUserFindByIdRes = "UserFindByIdRes"
	TypeUserFindAll     = "UserFindAll"
	TypeUserFindAllRes  = "UserFindAllRes"
	TypeUserUpdate      = "UserUpdate"

	TypeAirConditionerFindByRoom    = "AirConditionerFindByRoom"
	TypeAirConditionerFindByRoomRes = "AirConditionerFindByRoomRes"
	TypeAirConditionerFindAll       = "AirConditionerFindAll"
	TypeAirConditionerFindAllRes    = "AirConditionerFindAllRes"
	TypeAirConditionerCreate        = "AirConditionerCreate"
	TypeAirConditionerUpdate        = "AirConditionerUpdate"
	TypeAirConditionerOn            = "AirConditionerOn"
	TypeAirConditionerOff           = "AirConditionerOff"
	TypeAirConditionerSetParam      = "AirConditionerSetParam"
	TypeAirConditionerStopWind      = "AirConditionerStopWind"
	TypeGetReport                   = "GetReport"
	TypeGetReportRes                = "GetReportRes"
	TypeSetRoomData					= "SetRoomData"
	TypeGetDetailList				= "GetDetailList"
	TypeGetDetailListRes			= "GetDetailListRes"

	TypeFeeAdd       = "FeeAdd"
	TypeFeeQuery     = "FeeQuery"
	TypeFeeQueryRes  = "FeeQueryRes"
	TypeFeeDelete    = "FeeDelete"

	TypeGetServingQueue = "GetServingQueue"
	TypeGetServingQueueRes = "GetServingQueueRes"
)

// 定义消息结构体
type Message struct {
	Length int    `json:"length"` // 消息长度，用于验证包是否缺失
	Type   string `json:"type"`   // 消息类型
	Data   string `json:"data"`   // 消息
}

// 普通的消息返回格式，只包含状态码和msg，如需包含额外数据则另外定义返回格式
type NormalRes struct {
	Code int    `json:"code"` // 状态码
	Msg  string `json:"msg"`  // 返回信息
}

// 用户注册消息结构体
type UserRegister struct {
	RoomNum  int    `json:"room_num"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// 用户登录消息结构体
type UserLogin struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type UserFindById struct {
	Id int `json:"id"`
}

type UserFindByIdRes struct {
	Code int        `json:"code"` // 状态码
	Msg  string     `json:"msg"`  // 返回信息
	User model.User `json:"user"`
}

type UserFindAllRes struct {
	Code  int          `json:"code"` // 状态码
	Msg   string       `json:"msg"`  // 返回信息
	Users []model.User `json:"users"`
}

type UserUpdate struct {
	User model.User `json:"user"`
}

//
type AirConditionerFindById struct {
	Id int `json:"id"`
}

type AirConditionerFindByIdRes struct {
	Code           int                  `json:"code"`            // 状态码
	Msg            string               `json:"msg"`             // 返回信息
	AirConditioner model.AirConditioner `json:"air_conditioner"` // 返回一个空调状态数据
}

type AirConditionerFindByRoom struct {
	RoomNum int `json:"room_num"`
}

type AirConditionerFindByRoomRes struct {
	Code            int                  `json:"code"`            // 状态码
	Msg             string               `json:"msg"`             // 返回信息
	AirConditioner model.AirConditioner `json:"air_conditioner"` // 返回空调状态数据
}

type AirConditionerFindAllRes struct {
	Code            int                    `json:"code"`             // 状态码
	Msg             string                 `json:"msg"`              // 返回信息
	AirConditioners []model.AirConditioner `json:"air_conditioners"` // 返回一组空调状态数据
}

type AirConditionerCreate struct {
	AirConditioner model.AirConditioner `json:"air_conditioner"`
}

type AirConditionerUpdate struct {
	AirConditioner model.AirConditioner `json:"air_conditioner"`
}

type AirConditionerOn struct {
	RoomNum     int     `json:"room_num"`    // 房间号
	Mode        string     `json:"mode"`        // 模式
	WindLevel   string     `json:"wind_level"`  // 风速
	Temperature float64 `json:"temperature"` // 目标温度
	OpenTime    int64   `json:"open_time"`   // 开机时间，时间戳格式
}

type AirConditionerOff struct {
	RoomNum   int   `json:"room_num"`
	CloseTime int64 `json:"close_time"` // 关机时间，时间戳格式
}

type AirConditionerSetParam struct {
	RoomNum     int     `json:"room_num"`    // 房间号
	Mode        string     `json:"mode"`        // 模式
	WindLevel   string     `json:"wind_level"`  // 风速
	Temperature float64 `json:"temperature"` // 目标温度
}

type AirConditionerStopWind struct {
	RoomNum      int   `json:"room_num"`
	StopWindTime int64 `json:"stop_wind_time"` // 停止送风时间
}

type SetRoomData struct {
	RoomNum      int   	 `json:"room_num"`
	TotalPower 	 float64 `json:"total_power"` // 总耗电量
	RoomTemperature float64 `json:"room_temperature"` // 房间温度
}

//
type FeeAdd struct {
	model.Fee `json:"fee"`
}

type FeeQuery struct {
	RoomNum int `json:"room_num"`
}

type FeeQueryRes struct {
	Code int       `json:"code"` // 状态码
	Msg  string    `json:"msg"`  // 返回信息
	Fee  model.Fee `json:"fee"`
}

type FeeDelete struct {
	RoomNum   int   `json:"room_num"`
}

type GetReportRes struct {
	Code    int            `json:"code"`    // 状态码
	Msg     string         `json:"msg"`     // 返回信息
	Reports []model.Report `json:"reports"` // 报表数组
}

type GetDetailList struct {
	RoomNum      int   	 `json:"room_num"`
}


type GetDetailListRes struct {
	Code    int            `json:"code"`    // 状态码
	Msg     string         `json:"msg"`     // 返回信息
	Detail  model.Detail	`json:"detail"` // 详单
}

type GetServingQueueRes struct {
	Code    int            `json:"code"`    		// 状态码
	Msg     string         `json:"msg"`     		// 返回信息
	ServingQueue  []int		`json:"serving_queue"` 	// 服务队列
}