package processor

import (
	"encoding/json"
	"fmt"
	"github.com/wxmsummer/AirConditioner/common/message"
	"github.com/wxmsummer/AirConditioner/common/utils"
	"github.com/wxmsummer/AirConditioner/server/model"
	"github.com/wxmsummer/AirConditioner/server/repository"
	"net"
)

type AdminProcessor struct {
	Conn net.Conn
	Orm  *repository.AdminOrm
}

func (adp *AdminProcessor) AdminSignUp(msg *message.Message) (err error) {
	signUpRes 	:= message.NormalRes {}
	req 		:= message.AdminRegister {}

	err = json.Unmarshal([]byte(msg.Data), &req)
	if err != nil {
		fmt.Println(err)
	}

	_, err = adp.Orm.FindByField("user_name", req.UserName, "user_name")
	if err == nil {
		signUpRes.Code = -1
		signUpRes.Msg = "Username already exists!"
	} else {
		adm := model.Admin {
			UserName:      	req.UserName,
			Password:       req.Password,
			AuthorityLevel: req.AuthorityLevel,
		}
		err := adp.Orm.Create(&adm)
		if err != nil {
			fmt.Println(err)
			signUpRes.Code = -2
			signUpRes.Msg = "Failed to create user!"
		} else {
			signUpRes.Code = 0
			signUpRes.Msg = "Succeed to create user!"
		}
	}

	data, err := json.Marshal(signUpRes)
	if err != nil {
		fmt.Println(err)
	}

	var resp message.Message
	resp.Type = message.TypeAdminRegister
	resp.Data = string(data)

	data, err = json.Marshal(resp)

	tf := &utils.Transfer{Conn: adp.Conn}
	err = tf.WritePkg(data)

	return
}

func (adp *AdminProcessor) AdminSignIn(msg *message.Message) (err error) {
	signUpRes 	:= message.NormalRes {}
	req 		:= message.AdminLogin {}

	err = json.Unmarshal([]byte(msg.Data), &req)
	if err != nil {
		fmt.Println(err)
	}

	adm, err := adp.Orm.FindByField("user_name", req.UserName, "authority_level, password")
	if err != nil {
		signUpRes.Code = -1
		signUpRes.Msg = "User do not exists!"
	} else {
		if adm.Password != req.Password {
			signUpRes.Code = -2
			signUpRes.Msg = "Wrong password!"
		} else if adm.AuthorityLevel != req.AuthorityLevel {
			signUpRes.Code = -3
			signUpRes.Msg = "Current authority level does match target level"
		} else {
			signUpRes.Code = 0
			signUpRes.Msg = "Succeed to login"
		}
	}

	data, err := json.Marshal(signUpRes)
	if err != nil {
		fmt.Println(err)
	}

	var resp message.Message
	resp.Type = message.TypeAdminLogin
	resp.Data = string(data)

	data, _ = json.Marshal(resp)

	tf := &utils.Transfer{Conn: adp.Conn}
	err = tf.WritePkg(data)

	return err
}