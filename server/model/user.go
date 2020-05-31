package model

// 用户结构
type User struct {
	Id       int    `json:"id"`
	RoomNum  int    `json:"room_num"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	CheckIn  int64  `json:"check_in"`  //入住时间
	CheckOut int64  `json:"check_out"` //退房时间
}
