package model

import "github.com/pkg/errors"

// 根据业务逻辑需要，自定义一些错误
var (
	ErrorUserNotExists = errors.New("用户不存在，请注册后登录")
	ErrorUserExists    = errors.New("用户已存在，请重新注册")
	ErrorUserPwdWrong  = errors.New("用户密码不正确，请重新输入密码")
)
