package model

import (
	"AirConditioner/server/db"
	"fmt"
	"strconv"
)

// 用户结构
type User struct {
	RoomNum   int    `json:"room_num"`
	Privilege int    `json:"privilege"`
	Password  string `json:"password"`
	CheckIn   int64  `json:"check_in"`  //入住时间
	CheckOut  int64  `json:"check_out"` //退房时间
}

// 添加用户
func AddUser(user *User) (int64, error) {
	sql := "insert into users(roomNum, privilege, password, checkIn, checkOut) values (?, ?, ?, ?, ?)"
	return db.ModifyDB(sql, user.RoomNum, user.Privilege, user.Password, user.CheckIn, user.CheckOut)
}

// 通过房间号查询用户
func QueryUserWithRoomNum(roomNum int) User {
	sql := "select from users where roomNum = ?" + strconv.Itoa(roomNum)
	row := db.QueryRowDB(sql)
	user := User{}
	row.Scan(&user.RoomNum, &user.Privilege, &user.Password, &user.CheckIn, &user.CheckOut)
	return user
}

// 更新用户的入住时间
func UpdateUserCheckIn(roomNum int, checkIn int64) (int64, error) {
	sql := "update users set checkIn = ? where roomNum = ?"
	return db.ModifyDB(sql, checkIn, roomNum)
}

// 更新用户的退房时间
func UpdateUserCheckOut(roomNum int, checkOut int64) (int64, error) {
	sql := "update users set checkOut = ? where roomNum = ?"
	return db.ModifyDB(sql, checkOut, roomNum)
}

// 查询所有用户
func QueryAllUsers() (users []User, err error) {
	sql := "select * from users"
	rows, err := db.QueryRowsDB(sql)
	if err != nil {
		fmt.Println("QueryRowsDB err=", err)
		return nil, err
	}
	for rows.Next() {
		user := User{}
		rows.Scan(user.RoomNum, user.Privilege, user.Password, user.CheckIn, user.CheckOut)
		users = append(users, user)
	}
	return users, nil
}
