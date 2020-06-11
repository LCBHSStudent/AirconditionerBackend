package processor

import (
	"encoding/json"
	"fmt"
	"github.com/wxmsummer/AirConditioner/common/message"
	"github.com/wxmsummer/AirConditioner/common/utils"
	"github.com/wxmsummer/AirConditioner/server/model"
	"github.com/wxmsummer/AirConditioner/server/repository"
	"net"
	"time"
)

type UserProcessor struct {
	Conn net.Conn
	Orm  *repository.UserOrm
}

// 用户注册
func (up *UserProcessor) Register(msg *message.Message) (err error) {
	var userRegister message.UserRegister
	err = json.Unmarshal([]byte(msg.Data), &userRegister)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}

	var resMsg message.Message
	var userRegisterRes message.NormalRes

	isExistUser, err := up.Orm.FindByPhone(userRegister.Phone)
	if err != nil {
		fmt.Println("up.Orm.FindByPhone(userRegister.Phone) err=", err)
	}
	if isExistUser != nil {
		userRegisterRes.Code = 403
		userRegisterRes.Msg = message.ErrorUserExists
		return
	}

	user := &model.User{
		RoomNum:  userRegister.RoomNum,
		Phone:    userRegister.Phone,
		Password: userRegister.Password,
		CheckIn:  time.Now().Unix(),
		CheckOut: 0,
	}
	err = up.Orm.Create(user)
	if err != nil {
		fmt.Println("up.Orm.Create(user) err = ", err)
		return
	}

	userRegisterRes.Code = 200
	userRegisterRes.Msg = "注册用户成功！"

	data, err := json.Marshal(userRegisterRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeNormalRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: up.Conn}
	err = tf.WritePkg(data)
	return
}

// 用户登录
func (up *UserProcessor) Login(msg *message.Message) (err error) {
	var userLogin message.UserLogin
	err = json.Unmarshal([]byte(msg.Data), &userLogin)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}

	var resMsg message.Message
	var userLoginRes message.NormalRes

	phone := userLogin.Phone
	password := userLogin.Password

	existUser, err := up.Orm.FindByPhone(phone)
	if err != nil {
		fmt.Println("up.Orm.FindByPhone(userLogin.Phone) err =", err)
		return
	}

	if existUser.Id == 0 {
		userLoginRes.Code = 404 // 用户不存在，找不到资源
		userLoginRes.Msg = message.ErrorUserNotExists
		return
	}
	if existUser.Password != password {
		userLoginRes.Code = 401
		userLoginRes.Msg = message.ErrorUserPwdWrong
		return
	}

	userLoginRes.Code = 200 // 登陆成功
	userLoginRes.Msg = "用户登陆成功!"

	data, err := json.Marshal(userLoginRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeNormalRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: up.Conn}
	err = tf.WritePkg(data)
	return
}

func (up *UserProcessor) FindById(msg *message.Message) (err error) {
	var userFindById message.UserFindById
	err = json.Unmarshal([]byte(msg.Data), &userFindById)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}

	var resMsg message.Message
	var userFindByIdRes message.UserFindByIdRes

	id := userFindById.Id
	user, err := up.Orm.FindById(id)
	if err != nil {
		fmt.Println("up.Orm.FindByID(id) err=", err)
		return err
	}

	userFindByIdRes.Code = 200
	userFindByIdRes.Msg = "查询成功！"
	userFindByIdRes.User = user

	data, err := json.Marshal(userFindByIdRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeUserFindByIdRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: up.Conn}
	err = tf.WritePkg(data)
	return
}

func (up *UserProcessor) FindAll(msg *message.Message) (err error) {
	var resMsg message.Message
	var userFindAllRes message.UserFindAllRes

	users, err := up.Orm.FindAllUsers()
	if err != nil {
		fmt.Println("up.Orm.FindAllUsers() err=", err)
		return err
	}

	userFindAllRes.Code = 200
	userFindAllRes.Msg = "查询成功！"
	userFindAllRes.Users = users

	data, err := json.Marshal(userFindAllRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeUserFindAllRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: up.Conn}
	err = tf.WritePkg(data)
	return
}

func (up *UserProcessor) Update(msg *message.Message) (err error) {
	var userUpdate message.UserUpdate
	err = json.Unmarshal([]byte(msg.Data), &userUpdate)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}

	var resMsg message.Message
	var userUpdateRes message.NormalRes

	user := userUpdate.User
	err = up.Orm.Update(&user)
	if err != nil {
		fmt.Println("up.Orm.Update(&user) err=", err)
		return err
	}

	userUpdateRes.Code = 200
	userUpdateRes.Msg = "更新用户信息成功！"

	data, err := json.Marshal(userUpdateRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeNormalRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: up.Conn}
	err = tf.WritePkg(data)
	return
}
