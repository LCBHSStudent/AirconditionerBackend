package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/wxmsummer/AirConditioner/server/model"
)

// 空调数据库操作相关接口
type AirConditionerRepository interface {
	FindByID(int) (*model.AirConditioner, error)                       // 根据id查找空调
	Create(*model.AirConditioner) error                                // 创建空调状态
	Update(*model.AirConditioner) error                                // 更新空调状态
	FindByField(string, string, string) (*model.AirConditioner, error) // 条件查询
	FindAll() ([]model.AirConditioner, error)                          // 查询所有空调
	FindByRoom(int) ([]model.AirConditioner, error)                    // 通过房间号查询
}

type AirConditionerOrm struct {
	Db *gorm.DB
}

func (orm *AirConditionerOrm) FindByID(id int) (airConditioner model.AirConditioner, err error) {
	err = orm.Db.Where("id = ?", id).First(airConditioner).Error
	if err != nil {
		return
	}
	return
}

func (orm *AirConditionerOrm) Create(airConditioner *model.AirConditioner) error {
	err := orm.Db.Create(airConditioner).Error
	if err != nil {
		return err
	}
	return nil
}

func (orm *AirConditionerOrm) Update(airConditioner model.AirConditioner) error {
	err := orm.Db.Model(airConditioner).Update(&airConditioner).Error
	if err != nil {
		return err
	}
	return nil
}

func (orm *AirConditionerOrm) FindByField(key, value, fields string) (*model.AirConditioner, error) {
	if len(fields) == 0 {
		fields = "*"
	}
	airConditioner := &model.AirConditioner{}
	err := orm.Db.Select(fields).Where(key+" = ?", value).First(airConditioner).Error
	if err != nil {
		return nil, err
	}
	return airConditioner, nil
}

// 查询所有空调
func (orm *AirConditionerOrm) FindAll() (airs []model.AirConditioner, err error) {
	err = orm.Db.Find(&airs).Error
	if err != nil {
		return nil, err
	}
	return airs, nil
}

func (orm *AirConditionerOrm) FindByRoom(roomNum int) (airs []model.AirConditioner, err error) {
	err = orm.Db.Where("room_num = ?", roomNum).Find(&airs).Error
	if err != nil {
		return nil, err
	}
	return airs, nil
}
