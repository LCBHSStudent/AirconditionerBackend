package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/wxmsummer/AirConditioner/server/model"
)

// 账单数据库操作相关接口
type FeeRepository interface {
	AddFee(roomState *model.Fee) error                                    // 添加一条账单信息
	QueryFees(roomNum int, startTime, endTime int64) ([]model.Fee, error) // 查询指定时间段内的账单，返回账单数组
	DelFees(roomNum int, startTime, endTime int64) (int64, error)         // 删除某个房间指定时间段内的账单记录（删除太老的数据），返回删除条数
}

type FeeOrm struct {
	Db *gorm.DB
}

// 添加一条账单信息到数据库
func (orm *FeeOrm) AddFee(roomState *model.Fee) error {
	err := orm.Db.Create(roomState).Error
	if err != nil {
		return err
	}
	return nil
}

// 查询指定时间段内的账单信息
func (orm *FeeOrm) QueryFees(roomNum int, startTime, endTime int64) (roomStates []model.Fee, err error) {
	err = orm.Db.Where("roomNum = ? AND startTime >= ? AND endTime <= ?", roomNum, startTime, endTime).Find(&roomStates).Error
	if err != nil {
		return nil, err
	}
	return roomStates, nil
}

// 删除某个房间指定时间段内的账单记录（删除太老的数据），返回删除条数
func (orm *FeeOrm) DelFees(roomNum int, startTime, endTime int64) (int64, error) {
	roomState := &model.Fee{}
	dbResult := orm.Db.Where("roomNum = ? AND startTime >= ? AND endTime <= ?", roomNum, startTime, endTime).Delete(&roomState)
	err := dbResult.Error
	if err != nil {
		return 0, err
	}
	return dbResult.RowsAffected, nil
}
