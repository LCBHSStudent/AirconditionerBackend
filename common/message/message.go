package message

import "github.com/wxmsummer/AirConditioner/server/model"

// 消息类型
const (
	TypeUserRegister    = "UserRegister"
	TypeUserRegisterRes = "UserRegisterRes"
	TypeUserLogin       = "UserLogin"
	TypeUserLoginRes    = "UserLoginRes"
	TypeUserFindById    = "UserFindById"
	TypeUserFindByIdRes = "UserFindByIdRes"
	TypeUserFindAll     = "UserFindAll"
	TypeUserFindAllRes  = "UserFindAllRes"
	TypeUserUpdate      = "UserUpdate"
	TypeUserUpdateRes   = "UserUpdate"

	TypeAirConditionerFindById      = "AirConditionerFindById"
	TypeAirConditionerFindByIdRes   = "AirConditionerFindByIdRes"
	TypeAirConditionerFindByRoom    = "AirConditionerFindByRoom"
	TypeAirConditionerFindByRoomRes = "AirConditionerFindByRoomRes"
	TypeAirConditionerFindAll       = "AirConditionerFindAll"
	TypeAirConditionerFindAllRes    = "AirConditionerFindAllRes"
	TypeAirConditionerCreate        = "AirConditionerCreate"
	TypeAirConditionerCreateRes     = "AirConditionerCreateRes"
	TypeAirConditionerUpdate        = "AirConditionerUpdate"
	TypeAirConditionerUpdateRes     = "AirConditionerUpdateRes"

	TypeRoomStateAdd       = "RoomStateAdd"
	TypeRoomStateAddRes    = "RoomStateAddRes"
	TypeRoomStateQuery     = "RoomStateQuery"
	TypeRoomStateQueryRes  = "RoomStateQueryRes"
	TypeRoomStateDelete    = "RoomStateDelete"
	TypeRoomStateDeleteRes = "RoomStateDeleteRes"

	TypeFeeAdd       = "FeeAdd"
	TypeFeeAddRes    = "FeeAddRes"
	TypeFeeQuery     = "FeeQuery"
	TypeFeeQueryRes  = "FeeQueryRes"
	TypeFeeDelete    = "FeeDelete"
	TypeFeeDeleteRes = "FeeDeleteRes"
)

// 定义消息结构体
type Message struct {
	Length int    `json:"length"` // 消息长度，用于验证包是否缺失
	Type   string `json:"type"`   // 消息类型
	Data   string `json:"data"`   // 消息
}

// 用户注册消息结构体
type UserRegister struct {
	RoomNum  int    `json:"room_num"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// 用户注册返回的消息
type UserRegisterRes struct {
	Code int    `json:"code"`  // 状态码 400 表示用户已被注册，200表示注册成功
	Msg  string `json:"error"` // 返回错误信息
}

// 用户登录消息结构体
type UserLogin struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// 用户登录返回的消息
type UserLoginRes struct {
	Code int    `json:"code"`  // 状态码
	Msg  string `json:"error"` // 返回错误信息
}

type UserFindById struct {
	Id int `json:"id"`
}

type UserFindByIdRes struct {
	Code int        `json:"code"`  // 状态码
	Msg  string     `json:"error"` // 返回错误信息
	User model.User `json:"user"`
}

type UserFindAll struct {
}

type UserFindAllRes struct {
	Code  int          `json:"code"`  // 状态码
	Msg   string       `json:"error"` // 返回错误信息
	Users []model.User `json:"users"`
}

type UserUpdate struct {
	User model.User `json:"user"`
}

type UserUpdateRes struct {
	Code int    `json:"code"`  // 状态码
	Msg  string `json:"error"` // 返回错误信息
}

//
//
//
type AirConditionerFindById struct {
	Id int `json:"id"`
}

type AirConditionerFindByIdRes struct {
	Code           int                  `json:"code"`            // 状态码
	Msg            string               `json:"error"`           // 返回错误信息
	AirConditioner model.AirConditioner `json:"air_conditioner"` // 返回一个空调状态数据
}

type AirConditionerFindByRoom struct {
	RoomNum int `json:"room_num"`
}

type AirConditionerFindByRoomRes struct {
	Code            int                    `json:"code"`             // 状态码
	Msg             string                 `json:"error"`            // 返回错误信息
	AirConditioners []model.AirConditioner `json:"air_conditioners"` // 返回一组空调状态数据
}

type AirConditionerFindAll struct{}

type AirConditionerFindAllRes struct {
	Code            int                    `json:"code"`             // 状态码
	Msg             string                 `json:"error"`            // 返回错误信息
	AirConditioners []model.AirConditioner `json:"air_conditioners"` // 返回一组空调状态数据
}

type AirConditionerCreate struct {
	AirConditioner model.AirConditioner `json:"air_conditioner"`
}

type AirConditionerCreateRes struct {
	Code int    `json:"code"`  // 状态码
	Msg  string `json:"error"` // 返回错误信息
}

type AirConditionerUpdate struct {
	AirConditioner model.AirConditioner `json:"air_conditioner"`
}

type AirConditionerUpdateRes struct {
	Code int    `json:"code"`  // 状态码
	Msg  string `json:"error"` // 返回错误信息
}

//
//
//
type RoomStateAdd struct {
	model.RoomState `json:"room_state"`
}

type RoomStateAddRes struct {
	Code int    `json:"code"`  // 状态码
	Msg  string `json:"error"` // 返回错误信息
}

type RoomStateQuery struct {
	RoomNum   int   `json:"room_num"`
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

type RoomStateQueryRes struct {
	Code       int               `json:"code"`  // 状态码
	Msg        string            `json:"error"` // 返回错误信息
	RoomStates []model.RoomState `json:"room_states"`
}

type RoomStateDelete struct {
	RoomNum   int   `json:"room_num"`
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

type RoomStateDeleteRes struct {
	Code int    `json:"code"`  // 状态码
	Msg  string `json:"error"` // 返回错误信息
}

//
//
//
type FeeAdd struct {
	model.Fee `json:"fee"`
}

type FeeAddRes struct {
	Code int    `json:"code"`  // 状态码
	Msg  string `json:"error"` // 返回错误信息
}

type FeeQuery struct {
	RoomNum   int   `json:"room_num"`
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

type FeeQueryRes struct {
	Code int         `json:"code"`  // 状态码
	Msg  string      `json:"error"` // 返回错误信息
	Fees []model.Fee `json:"fees"`
}

type FeeDelete struct {
	RoomNum   int   `json:"room_num"`
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

type FeeDeleteRes struct {
	Code int    `json:"code"`  // 状态码
	Msg  string `json:"error"` // 返回错误信息
}
