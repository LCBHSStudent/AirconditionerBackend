package Processor

import (
	"encoding/json"
	"fmt"
	"github.com/wxmsummer/airConditioner/server/message"
	"github.com/wxmsummer/airConditioner/server/model"
	"github.com/wxmsummer/airConditioner/server/utils"
	"net"
)

type UserProcessor struct {
	Conn    net.Conn
	RoomNum int
}

// 用户注册
func (this *UserProcessor) ProcessUserRegister(msg *message.Message) (err error) {
	var userRegister message.UserRegister
	err = json.Unmarshal([]byte(msg.Data), &userRegister)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}
	var resMsg message.Message
	var userRegisterRes message.UserRegisterRes

	user1 := model.QueryUserWithRoomNum(userRegister.RoomNum)
	if user1.RoomNum != 0 { // 说明该用户已经注册
		userRegisterRes.Code = 403
		userRegisterRes.Error = model.ErrorUserExists.Error()
	} else {
		user := &model.User{
			RoomNum:   userRegister.RoomNum,
			Privilege: userRegister.Privilege,
			Password:  userRegister.Password,
			CheckIn:   0,
			CheckOut:  0,
		}
		_, err = model.AddUser(user)
		if err != nil {
			fmt.Println("db add user err = ", err)
			userRegisterRes.Code = 500
			userRegisterRes.Error = err.Error()
		} else {
			fmt.Println("添加用户成功...")
			userRegisterRes.Code = 200
		}
	}
	data, err := json.Marshal(userRegisterRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}
	resMsg.Type = message.TypeUserRegisterRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}
	tf := &utils.Transfer{Conn: this.Conn}
	err = tf.WritePkg(data)
	return
}

// 用户登录
func (this *UserProcessor) ProcessUserLogin(msg *message.Message) (err error) {
	var userLogin message.UserLogin
	err = json.Unmarshal([]byte(msg.Data), &userLogin)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}
	var resMsg message.Message
	var userLoginRes message.UserLoginRes
	roomNum := userLogin.RoomNum
	password := userLogin.Password
	user := model.QueryUserWithRoomNum(roomNum)
	fmt.Println("query user:", user)
	if user.RoomNum == 0 {
		err = model.ErrorUserNotExists
		userLoginRes.Code = 404 // 用户不存在，找不到资源
		userLoginRes.Error = err.Error()
		return
	} else if user.Password != password {
		err = model.ErrorUserPwdWrong
		userLoginRes.Code = 401 // 密码错误，请求重新认证
		userLoginRes.Error = err.Error()
		return
	} else if user.RoomNum == roomNum && user.Password == password {
		userLoginRes.Code = 200 // 登陆成功
		fmt.Println("用户登陆成功...")
	}
	data, err := json.Marshal(userLoginRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}
	resMsg.Type = message.TypeUserLoginRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}
	tf := &utils.Transfer{Conn: this.Conn}
	err = tf.WritePkg(data)
	return
}

// 修改用户的入住时间或者退房时间
func (this *UserProcessor) ProcessUserModify(msg *message.Message) (err error) {
	var userModify message.UserModify
	err = json.Unmarshal([]byte(msg.Data), &userModify)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}
	var resMsg message.Message
	var userModifyRes message.UserModifyRes
	roomNum := userModify.RoomNum
	checkIn := userModify.CheckIn
	checkOut := userModify.CheckOut
	if checkIn != 0 && checkOut == 0 { //说明是用户入住
		_, err = model.UpdateUserCheckIn(roomNum, checkIn)
		if err == nil {
			userModifyRes.Code = 200 // 修改入住时间成功
			fmt.Println("修改用户入住时间成功...")
		} else {
			userModifyRes.Code = 500
			userModifyRes.Error = err.Error()
		}
	} else if checkIn == 0 && checkOut != 0 { // 说明是用户退房
		// 1，到数据库修改用户的退房时间
		// 2，从数据库查询账单、详单信息，返回
		_, err = model.UpdateUserCheckOut(roomNum, checkOut)
		if err == nil {
			userModifyRes.Code = 200 // 修改退房时间成功
			fmt.Println("修改用户退房时间成功...")
			user := model.QueryUserWithRoomNum(roomNum) // 查询用户，获取他的入住时间
			roomStates, err := model.QueryRoomStates(roomNum, user.CheckIn, checkOut)
			if err == nil { // 说明查询到了详单数据
				userModifyRes.Code = 200
				userModifyRes.RoomStateList = roomStates
				fmt.Println("查询用户住房详单成功...")
			} else { // 否则报错
				userModifyRes.Code = 500
				userModifyRes.Error = err.Error()
			}
		} else {
			userModifyRes.Code = 500
			userModifyRes.Error = err.Error()
		}
	}
	data, err := json.Marshal(userModifyRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}
	resMsg.Type = message.TypeUserModifyRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}
	tf := &utils.Transfer{Conn: this.Conn}
	err = tf.WritePkg(data)
	return
}

// 用户查询 如果传入的房间号为0，则查询所有用户信息，否则按房间号查询
func (this *UserProcessor) ProcessUserQuery(msg *message.Message) (err error) {
	var userQuery message.UserQuery
	err = json.Unmarshal([]byte(msg.Data), &userQuery)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}
	var resMsg message.Message
	var userQueryRes message.UserQueryRes
	roomNum := userQuery.RoomNumber
	if roomNum != 0 { // 查询单个用户的信息
		user := model.QueryUserWithRoomNum(roomNum)
		var userList []model.User
		userList = append(userList, user)
		userQueryRes.Code = 200
		userQueryRes.UserList = userList
	} else if roomNum == 0 { // 查询所有用户的信息
		userList, err := model.QueryAllUsers()
		if err == nil {
			userQueryRes.Code = 200
			userQueryRes.UserList = userList
		} else {
			userQueryRes.Code = 500
			userQueryRes.Error = err.Error()
		}
	}
	data, err := json.Marshal(userQueryRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}
	resMsg.Type = message.TypeUserQueryRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}
	tf := &utils.Transfer{Conn: this.Conn}
	err = tf.WritePkg(data)
	return
}
