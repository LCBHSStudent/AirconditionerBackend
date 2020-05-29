package model

import (
	"fmt"
	"github.com/wxmsummer/airConditioner/server/db"
	"strconv"
)

const (
	// 电源开关
	PowerOff = 0
	PowerOn  = 1

	// 空调模式
	ModeCold      = 0
	ModeHot       = 1
	ModeWind      = 2
	ModeDry       = 3
	ModeSleep     = 4
	ModeSwingFlap = 5
	ModeBreath    = 6

	// 风速
	WindAuto = 0
	WindLow  = 1
	WindMid  = 2
	WindHigh = 3
)

// 空调数据结构
type AirConditioner struct {
	Number      int     `json:"number"`      // 空调编号，默认和房间编号相同
	Power       int     `json:"power"`       // 电源开关：0关 1开
	Mode        int     `json:"mode"`        // 模式
	WindLevel   int     `json:"wind_level"`  // 风速
	Temperature float64 `json:"temperature"` // 温度
}

// 往数据库中增加一条空调记录
func AddAirConditioner(air *AirConditioner) (int64, error) {
	sql := "insert into airs(number, power, mode, windLever, temperature)values(?,?,?,?,?) "
	return db.ModifyDB(sql, air.Number, air.Power, air.Mode, air.WindLevel, air.Temperature)
}

// 通过空调编号查询空调状态
func QueryAirWithNumber(number int) (air AirConditioner) {
	sql := "select power, mode, windLevel, temperature from airs where number = ?" + strconv.Itoa(number)
	row := db.QueryRowDB(sql)
	row.Scan(&air.Power, &air.Mode, &air.WindLevel, &air.Temperature)
	return air
}

// 修改空调状态数据
func UpdateAirConditioner(air *AirConditioner) (int64, error) {
	sql := "update airs set power=?, mode=?, windLevel=?, temperature=? where number=?"
	return db.ModifyDB(sql, air.Power, air.Mode, air.WindLevel, air.Temperature, air.Number)
}

// 查询所有空调
func QueryAllAirConditioners() (airs []AirConditioner, err error) {
	sql := "select number, power, mode, windLevel, temperature from airs"
	rows, err := db.QueryRowsDB(sql)
	if err != nil {
		fmt.Println("QueryRowsDB err=", err)
		return nil, err
	}
	for rows.Next() {
		air := AirConditioner{}
		rows.Scan(air.Number, air.Power, air.Mode, air.WindLevel, air.Temperature)
		airs = append(airs, air)
	}
	return airs, nil
}
