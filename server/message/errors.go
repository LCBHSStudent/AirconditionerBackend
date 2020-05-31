package message

// 根据业务逻辑需要，自定义一些错误
const (
	ErrorUserNotExists = "用户不存在，请注册后登录"
	ErrorUserExists    = "用户已存在，请重新注册"
	ErrorUserPwdWrong  = "用户密码不正确，请重新输入密码"
)
