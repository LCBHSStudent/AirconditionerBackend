package repository

import (
	"github.com/wxmsummer/AirConditioner/server/model"
	"gorm.io/gorm"
)

// 空调数据库操作相关接口
type AdminRepository interface {
	Create(*model.Admin) error                                	// 创建管理员
	FindByField(string, string, string) (*model.Admin, error) 	// 条件查询
	FindAll() ([]model.Admin, error)                          	// 查询所有
}

type AdminOrm struct {
	Db *gorm.DB
}

func (orm *AdminOrm) Create(admin *model.Admin) error {
	dbResult := orm.Db.Create(admin)
	err := dbResult.Error
	if err != nil {
		return err
	}
	return nil
}

func (orm *AdminOrm) FindByField(key, value, fields string) (*model.Admin, error) {
	if len(fields) == 0 {
		fields = "*"
	}
	admin := &model.Admin{}
	err := orm.Db.Select(fields).Where(key+" = ?", value).First(admin).Error
	if err != nil {
		return nil, err
	}
	return admin, nil
}

// FindAll 查询所有管理员
func (orm *AdminOrm) FindAll() (admin []model.Admin, err error) {
	err = orm.Db.Find(&admin).Error
	if err != nil {
		return nil, err
	}
	return admin, nil
}
