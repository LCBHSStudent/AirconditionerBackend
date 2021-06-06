package processor

import (
	"encoding/json"
	"fmt"
	"github.com/wxmsummer/AirConditioner/common/message"
	"github.com/wxmsummer/AirConditioner/common/utils"
	"github.com/wxmsummer/AirConditioner/server/model"
	"github.com/wxmsummer/AirConditioner/server/repository"
	"log"
	"net"
	"strconv"
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

	feOrm := repository.FeeOrm{Db: up.Orm.Db}
	fee := model.Fee{
		RoomNum: userRegister.RoomNum,
		Cost:    0.0,
	}

	err = feOrm.Create(&fee)
	if err != nil {
		log.Println(err)
	}

	var resMsg message.Message
	var userRegisterRes message.NormalRes

	users, err := up.Orm.FindByField("phone", userRegister.Phone, "")
	if len(users) != 0 {
		userRegisterRes.Code = 403
		userRegisterRes.Msg = message.ErrorUserExists
	} else {
		users, err := up.Orm.FindByField("room_num", strconv.Itoa(userRegister.RoomNum), "")
		if len(users) != 0 {
			userRegisterRes.Code = 400
			userRegisterRes.Msg = message.ErrorRoomHasUser
			log.Println(err)
		}  else {
			airOrm := repository.AirConditionerOrm{Db: up.Orm.Db}
			room, err := airOrm.FindByRoom(userRegister.RoomNum)
			if  room.RoomNum == 0 || err != nil {
				userRegisterRes.Code = 400
				userRegisterRes.Msg = message.ErrorRoomNotExist
			} else {
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
					return err
				}

				userRegisterRes.Code = 200
				userRegisterRes.Msg = "客户登记成功！"
			}

		}
	}

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
	var userLoginRes message.UserLoginRes

	phone := userLogin.Phone
	password := userLogin.Password

	resMsg.Type = message.TypeNormalRes

	existUser, err := up.Orm.FindByField("phone", phone, "password, room_num")
	if len(existUser) == 0 {
		userLoginRes.Code = 404 // 用户不存在，找不到资源
		userLoginRes.Msg = message.ErrorUserNotExists
	} else {
		if existUser[0].Password != password {
			userLoginRes.Code = 401
			userLoginRes.Msg = message.ErrorUserPwdWrong
		} else {
			userLoginRes.Code = 0 // 登陆成功
			userLoginRes.Msg = "登陆成功!"
			userLoginRes.RoomNumber = existUser[0].RoomNum
			resMsg.Type = message.TypeUserLogin
		}
	}

	data, err := json.Marshal(userLoginRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

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
	var userFindById message.UserFindByRoom
	err = json.Unmarshal([]byte(msg.Data), &userFindById)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}

	var resMsg message.Message
	var userFindByIdRes message.UserFindByIdRes

	roomNumber := userFindById.RoomNumber
	user, err := up.Orm.FindById(roomNumber)
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

func (up *UserProcessor) Checkout (msg *message.Message) (err error) {
	var req message.UserCheckout
	err = json.Unmarshal([]byte(msg.Data), &req)
	if err != nil {
		return err
	}

	var resMsg message.Message
	resMsg.Type = message.TypeNormalRes

	var checkoutRes message.NormalRes

	users, err := up.Orm.FindByField("room_num",strconv.Itoa(req.RoomNum), "")
	if err != nil {
		checkoutRes.Code = 400
		checkoutRes.Msg = "查询用户失败"
	} else if len(users) == 0 {
		checkoutRes.Code = 400
		checkoutRes.Msg = "房间不存在或房间不存在用户"
	} else {
		up.Orm.Db.Model(&model.User{}).Delete(users[0])

		airOrm := repository.AirConditionerOrm{Db: up.Orm.Db}
		deInitAir := model.AirConditioner{
			RoomNum:         req.RoomNum,
			Power:           "",
			Mode:            "",
			WindLevel:       "",
			Temperature:     0,
			RoomTemperature: 0,
			TotalPower:      0,
			TotalFee:        0,
			StartWind:       "",
			StopWind:        "",
			OpenTime:        "",
			CloseTime:       "",
			SetParamNum:     0,
		}
		airOrm.Db.Delete(deInitAir, "room_num = " + strconv.Itoa(req.RoomNum))
		airOrm.Db.Create(deInitAir)

		fo := repository.FeeOrm{Db: up.Orm.Db}
		err = fo.Delete(req.RoomNum)
		if err != nil {
			return err
		}

		checkoutRes.Code = 200
		checkoutRes.Msg = "用户退房成功"
	}

	data, err := json.Marshal(checkoutRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

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