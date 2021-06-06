package repository

import (
	"github.com/wxmsummer/AirConditioner/server/model"
	"gorm.io/gorm"
)

// 账单数据库操作相关接口
type FeeRepository interface {
	Create(*model.Fee) error                // 创建一条账单信息
	QueryByRoomNum(int) (model.Fee, error) // 查询某房间的账单
	Delete(roomNum int) (error)         // 删除某个房间的账单记录
	Update(fee *model.Fee) error
}

type FeeOrm struct {
	Db *gorm.DB
}

func (orm *FeeOrm) Create(fee *model.Fee) error {
	err := orm.Db.Create(fee).Error
	if err != nil {
		return err
	}
	return nil
}

func (orm *FeeOrm) QueryByRoomNum(roomNum int) (fee model.Fee, err error) {
	err = orm.Db.Where("room_num = ? ", roomNum).First(&fee).Error
	if err != nil {
		return fee, err
	}
	return fee, nil
}

func (orm *FeeOrm) Delete(roomNum int) error {
	fee := &model.Fee{}
	err := orm.Db.Where("room_num = ?", roomNum).Delete(fee).Error
	if err != nil {
		return err
	}
	return nil
}

func (orm *FeeOrm) Update(fee *model.Fee) error {
	err := orm.Db.Where("room_num = ?", fee.RoomNum).Save(fee).Error
	if err != nil {
		return err
	}
	return nil
}