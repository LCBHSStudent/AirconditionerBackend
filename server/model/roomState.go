package model

import (
	"fmt"
	"github.com/wxmsummer/airConditioner/server/db"
)

// 房间状态结构体，记录房间号、开始时间、结束时间，以及这段时间内的温度和费用、耗电量
type RoomState struct {
	RoomNum     int     `json:"room_num"`    // 房间号
	StartTime   int64   `json:"start_time"`  // 开始时间
	EndTime     int64   `json:"end_time"`    // 结束时间
	Power       float64 `json:"power"`       // 耗电量
	Cost        float64 `json:"cost"`        // 费用
	Temperature float64 `json:"temperature"` // 温度
}

// 添加一条房间的温度、耗电量、费用记录
func AddRoomState(roomState *RoomState) (int64, error) {
	sql := "insert into roomStates(roomNum, startTime, endTime, power, cost, temperature) values (?,?,?,?,?,?)"
	return db.ModifyDB(sql, roomState.RoomNum, roomState.StartTime, roomState.EndTime, roomState.Power, roomState.Cost, roomState.Temperature)
}

// 查询指定时间段内的温度、耗电量、费用记录
func QueryRoomStates(roomNum int, startTime, endTime int64) (roomStates []RoomState, err error) {
	sql := fmt.Sprintf("select from roomStates where roomNum = %d and startTime >= %d and endTime <= %d", roomNum, startTime, endTime)
	rows, err := db.QueryRowsDB(sql)
	if err != nil {
		fmt.Println("QueryRowsDB err=", err)
		return nil, err
	}
	for rows.Next() {
		roomState := RoomState{}
		rows.Scan(&roomState.RoomNum, &roomState.StartTime, &roomState.EndTime, roomState.Power, &roomState.Cost, &roomState.Temperature)
		roomStates = append(roomStates, roomState)
	}
	return
}

// 删除某个房间指定时间段内的温度、耗电量、费用记录（删除太老的数据）
func DelRoomStatesByRoomNum(roomNum int, startTime, endTime int64) (int64, error) {
	sql := fmt.Sprintf("delete from roomStates where roomNum = %d and startTime >= %d and endTime <= %d", roomNum, startTime, endTime)
	return db.ModifyDB(sql)
}
