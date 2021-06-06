package message

// 根据业务逻辑需要，自定义一些错误
const (
	ErrorUserNotExists 	= "客户信息不存在，请入住登录"
	ErrorUserExists    	= "客户信息已存在，请检查数据"
	ErrorUserPwdWrong  	= "用户密码错误，请重新输入密码"
	ErrorRoomHasUser	= "房间已被入住，请更换房间号"
	ErrorRoomNotExist	= "房间号不存在，请重新输入"
)
