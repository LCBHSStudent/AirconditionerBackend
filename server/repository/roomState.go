package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/wxmsummer/AirConditioner/server/model"
)

// 房间状态数据库操作相关接口
type RoomStateRepository interface {
	AddRoomState(roomState *model.RoomState) error                                    // 添加一条房间状态信息
	QueryRoomStates(roomNum int, startTime, endTime int64) ([]model.RoomState, error) // 查询指定时间段内的房间状态，返回房间状态数组
	DelRoomStates(roomNum int, startTime, endTime int64) (int64, error)               // 删除某个房间指定时间段内的状态记录（删除太老的数据），返回删除条数
}

type RoomStateOrm struct {
	Db *gorm.DB
}

// 添加一条房间状态信息到数据库
func (orm *RoomStateOrm) AddRoomState(roomState *model.RoomState) error {
	err := orm.Db.Create(roomState).Error
	if err != nil {
		return err
	}
	return nil
}

// 查询指定时间段内的房间状态信息
func (orm *RoomStateOrm) QueryRoomStates(roomNum int, startTime, endTime int64) (roomStates []model.RoomState, err error) {
	err = orm.Db.Where("roomNum = ? AND startTime >= ? AND endTime <= ?", roomNum, startTime, endTime).Find(&roomStates).Error
	if err != nil {
		return nil, err
	}
	return roomStates, nil
}

// 删除某个房间指定时间段内的温度、耗电量、费用记录（删除太老的数据）
func (orm *RoomStateOrm) DelRoomStates(roomNum int, startTime, endTime int64) (int64, error) {
	roomState := &model.RoomState{}
	dbResult := orm.Db.Where("roomNum = ? AND startTime >= ? AND endTime <= ?", roomNum, startTime, endTime).Delete(&roomState)
	err := dbResult.Error
	if err != nil {
		return 0, err
	}
	return dbResult.RowsAffected, nil
}
