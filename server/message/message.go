package message

import "github.com/wxmsummer/airConditioner/server/model"

// 消息类型
const (
	TypeUserRegister    = "UserRegister"
	TypeUserRegisterRes = "UserRegisterRes"
	TypeUserLogin       = "UserLogin"
	TypeUserLoginRes    = "UserLoginRes"
	TypeUserModify      = "UserModify"
	TypeUserModifyRes   = "UserModifyRes"
	TypeUserQuery       = "UserQuery"
	TypeUserQueryRes    = "UserQueryRes"

	TypeAirModify    = "AirModify"
	TypeAirModifyRes = "AirModifyRes"
	TypeAirQuery     = "AirQuery"
	TypeAirQueryRes  = "AirQueryRes"

	TypeRoomStateAdd       = "RoomStateAdd"
	TypeRoomStateAddRes    = "RoomStateAddRes"
	TypeRoomStateQuery     = "RoomStateQuery"
	TypeRoomStateQueryRes  = "RoomStateQueryRes"
	TypeRoomStateDelete    = "RoomStateDelete"
	TypeRoomStateDeleteRes = "RoomStateDeleteRes"
)

// 定义消息结构体
type Message struct {
	Length int    `json:"length"` // 消息长度，用于验证包是否缺失
	Type   string `json:"type"`   // 消息类型
	Data   string `json:"data"`   // 消息
}

// 用户注册消息结构体
type UserRegister struct {
	RoomNum   int    `json:"room_num"`
	Privilege int    `json:"privilege"`
	Password  string `json:"password"`
}

// 用户注册返回的消息
type UserRegisterRes struct {
	Code  int    `json:"code"`  // 状态码 400 表示用户已被注册，200表示注册成功
	Error string `json:"error"` // 返回错误信息
}

// 用户登录消息结构体
type UserLogin struct {
	RoomNum  int    `json:"room_num"`
	Password string `json:"password"`
}

// 用户登录返回的消息
type UserLoginRes struct {
	Code  int    `json:"code"`  // 状态码
	Error string `json:"error"` // 返回错误信息
}

// 修改用户状态 结构体
type UserModify struct {
	model.User `json:"user"`
}

// 修改用户状态返回的消息结构体
type UserModifyRes struct {
	Code          int               `json:"code"`  // 状态码
	Error         string            `json:"error"` // 返回错误信息
	RoomStateList []model.RoomState `json:"room_state_list"`
}

// 查询用户状态 结构体
type UserQuery struct {
	RoomNumber int `json:"room_number"` // 传入房间号进行查询，，如果传入0000号，则查询所有房间
}

// 查询用户状态返回的消息结构体
type UserQueryRes struct {
	Code     int          `json:"code"`      // 状态码
	Error    string       `json:"error"`     // 返回错误信息
	UserList []model.User `json:"user_list"` // 返回一个用户列表
}

// 修改空调状态
type AirModify struct {
	model.AirConditioner `json:"air_conditioner"`
}

// 修改空调状态返回的消息
type AirModifyRes struct {
	Code  int    `json:"code"`  // 状态码
	Error string `json:"error"` // 返回错误信息
}

// 查询空调状态
type AirQuery struct {
	RoomNumber int `json:"room_number"` // 传入房间号进行查询，，如果传入0号，则查询所有房间
}

// 查询空调状态返回的消息结构体
type AirQueryRes struct {
	Code    int                    `json:"code"`     // 状态码
	Error   string                 `json:"error"`    // 返回错误信息
	AirList []model.AirConditioner `json:"air_list"` // 返回一个空调列表
}

// 添加房间状态信息
type RoomStateAdd struct {
	model.RoomState `json:"room_state"`
}

// 添加房间状态 返回消息
type RoomStateAddRes struct {
	Code  int    `json:"code"`  // 状态码
	Error string `json:"error"` // 返回错误信息
}

// 查询房间状态:温度、耗电量、费用
type RoomStateQuery struct {
	RoomNum   int   `json:"room_num"`
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

// 查询房间状态返回的消息结构体
type RoomStateQueryRes struct {
	Code          int               `json:"code"`  // 状态码
	Error         string            `json:"error"` // 返回错误信息
	RoomStateList []model.RoomState `json:"room_state_list"`
}

// 删除房间状态
type RoomStateDelete struct {
	RoomNum   int   `json:"room_num"`
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

// 删除房间状态返回的消息结构体
type RoomStateDeleteRes struct {
	Code  int    `json:"code"`  // 状态码
	Error string `json:"error"` // 返回错误信息
}
